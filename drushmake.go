package main

import (
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"regexp"
	"strings"
)

// Parse a makefile on filePath and return the list of modules, themes and core.
func parseMakefile(filePath string) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
	}

	// keys := cfg.Section("").Keys()
	// names := cfg.Section("").KeyStrings()
	hash := cfg.Section("").KeysHash()
	var componentList []string

	for key, val := range hash {
		// fmt.Println(keyMapper(key))
		if componentName := keyMapper(key); componentName != "" && val != "core" {
			componentList = append(componentList, componentName)
			// fmt.Printf("%v => %v\n", componentName, val)
		}
	}

	// Grab core information
	// rawCoreVersion := cfg.Section("").Key("core")
	// core := Component{CORE, "drupal", parseVersion(rawCoreVersion.Value(), 7)}

	// version := SemVersion{}
	// version.init(rawCoreVersion.Value())
	// core := Component{CORE, "drupal", version}

	// fmt.Printf("%v\n", core.printVersion())
}

func keyMapper(key string) string {
	match, start, end := strings.Index(key, "projects["), strings.Index(key, "["), strings.Index(key, "]")
	if match == 0 && start > 0 && end > 0 {
		return key[start+1 : end]
	} else {
		return ""
	}
}

// Parse a string block which belongs to a single component and return the Component.
func (C *Component) blockToComponentParser(block string, coreVersion int) {
	var versionLocation = map[string]string{
		"BASIC":    `^projects\[(\w+)\]\s?=\s?(\S+)$`,                         // projects[views] = 3.14
		"VERSION":  `^projects\[(\w+)\]\[version\]\s?=\s?(\S+)$`,              // projects[nodequeue][version] = 2.0-alpha1
		"BRANCH":   `^projects\[(\w+)\]\[download\]\[branch\]\s?=\s?(\S+)$`,   // projects[ns_core][download][branch] = 7.x-2.x
		"REVISION": `^projects\[(\w+)\]\[download\]\[revision\]\s?=\s?(\S+)$`, // projects[draggableviews][download][revision] = 9677bc18b7255e13c33ac3cca48732b855c6817d
	}

	componentName := keyMapper(block)
	var version string
	var revision string

	scanner := bufio.NewScanner(strings.NewReader(block))
	for scanner.Scan() {
		line := scanner.Text()

		for rowType, expression := range versionLocation {
			re := regexp.MustCompile(expression)
			if isMatch := re.MatchString(line); isMatch {
				matches := re.FindStringSubmatch(line)

				// Sanity check for entries not for a single project.
				if matches[1] != componentName {
					continue
				}

				// Do not overwrite the version variable when the revision value is captured.
				if rowType == "REVISION" {
					revision = matches[2]
				} else {
					version = matches[2]
				}

				// Make adjustments to the processed variables based on aggreageted information.
				if rowType == "BRANCH" {
					version = version + "-git"
				}

				fmt.Println(matches)
			}
		}
	}
	if revision != "" {
		version = version + "-git:" + revision
	}

	C.init(coreVersion, version, componentName, MODULE)

	// Structure variations:
	// 1. oneliner with version
	// 2. multiple lines with explicit version
	// 3. multiple lines with git
	// 3. multiple lines with git and hah
	// 5. multiple lines with dev version
}

// Helper function to determine the applicable structure.
func definitionClassifier(block string) {
	// project[component] = VERSION
	// project[component][download][type][git]
	// project[component][download][type][git][branch]
	// project[component][download][type][git][revision]
}

// Build a list of deduped project names out of a raw projects[name] block.
func componentList(rawBlock string) (componentList []string) {
	components := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(rawBlock))

	for scanner.Scan() {
		key := keyMapper(scanner.Text())

		if key == "" { // Skip empty lines.
			continue
		} else if _, no := components[key]; key != "" && !bool(no) {
			components[key] = true
		}
	}
	for name, _ := range components {
		componentList = append(componentList, name)
	}
	return
}

// Find a code block which contains a certain component key.
func findBlock(component string) (componentBlock string) {
	scanner := bufio.NewScanner(strings.NewReader(componentBlock))

	for scanner.Scan() {
		line := scanner.Text()
		// Find a line which contains our keyword.
		if match := strings.Index(line, fmt.Sprintf("projects[%s]", component)); match > -1 {
			component += line + "\n"
		}
	}
	return
}
