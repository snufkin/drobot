package main

import (
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"regexp"
	"strings"
)

// Parse a makefile and populate the manifest.
func (M *Manifest) parseMakefile(filePath string) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
	}
	// 0. Identify the core version and initialise the manifest.
	// 1. Identify the list of components.
	// 2. Extract the block relevant for each component.
	// 3. Process each block and populate the manifest.

	// 0. Extract the core information and init the manifest.
	rawCoreVersion := cfg.Section("").Key("core")
	core := Component{}
	// Core initialisation does not require an accurate core version, this will
	// be extracted from the raw version number.
	core.init(0, rawCoreVersion.Value(), "drupal", CORE)

	M.Components = append(M.Components, core)

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Print(err)
	}
	componentList := componentList(string(b))

	for _, cName := range componentList {
		c := Component{}
		c.blockToComponentParser(findBlock(cName, string(b)), core.Version.Major)
		M.Components = append(M.Components, c)
	}
}

// Helper function to extract the project name from a component block.
func keyMapper(key string) string {
	match, start, end := strings.Index(key, "projects["), strings.Index(key, "["), strings.Index(key, "]")
	if match == 0 && start > 0 && end > 0 {
		return key[start+1 : end]
	} else {
		return ""
	}
}

type blockInfo struct {
	Name     string
	Version  string
	Type     string
	Revision string
}

// Parse a string block which belongs to a single component and return the Component.
func (C *Component) blockToComponentParser(block string, coreVersion int) {
	var versionLocation = map[string]string{
		"BASIC":    `^projects\[(\w+)\]\s?=\s?(\S+)$`,                         // projects[views] = 3.14
		"VERSION":  `^projects\[(\w+)\]\[version\]\s?=\s?(\S+)$`,              // projects[nodequeue][version] = 2.0-alpha1
		"BRANCH":   `^projects\[(\w+)\]\[download\]\[branch\]\s?=\s?(\S+)$`,   // projects[ns_core][download][branch] = 7.x-2.x
		"TYPE":     `^projects\[(\w+)\]\[download\]\[type\]\s?=\s?(\S+)$`,     // projects[ns_core][download][type] = git
		"REVISION": `^projects\[(\w+)\]\[download\]\[revision\]\s?=\s?(\S+)$`, // projects[draggableviews][download][revision] = 9677bc18b7255e13c33ac3cca48732b855c6817d
	}

	componentName := keyMapper(block)
	component := blockInfo{Name: componentName}

	// We assume that a single block will reference a single component, see the continue.
	scanner := bufio.NewScanner(strings.NewReader(block))
	for scanner.Scan() {
		line := scanner.Text()

		for rowType, expression := range versionLocation {
			re := regexp.MustCompile(expression)
			if isMatch := re.MatchString(line); isMatch {
				matches := re.FindStringSubmatch(line)

				// Sanity check for entries not for a single project. All regex captures the name as match[1].
				if matches[1] != component.Name {
					continue
				}

				// Populate the right components within the struct.
				switch rowType {
				case "BASIC", "VERSION", "BRANCH":
					component.Version = matches[2]
				case "TYPE":
					component.Type = matches[2]
				case "REVISION":
					component.Revision = matches[2]
				}
			}
		}
	}

	C.init(coreVersion, component.String(), component.Name, MODULE)
}

func (b blockInfo) String() string {
	if b.Type != "" && b.Revision != "" { // Revision and Type set.
		return fmt.Sprintf("%s-%s:%s", b.Version, b.Type, b.Revision)
	} else if b.Type != "" && b.Revision == "" {
		return fmt.Sprintf("%s-%s", b.Version, b.Type)
	} else {
		return fmt.Sprintf("%s", b.Version)
	}
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
func findBlock(cName string, rawBlock string) (componentBlock string) {
	scanner := bufio.NewScanner(strings.NewReader(rawBlock))
	componentBlock = ""

	for scanner.Scan() {
		line := scanner.Text()
		// Find a line which contains our keyword.
		if match := strings.Index(line, fmt.Sprintf("projects[%s]", cName)); match > -1 {
			componentBlock += line + "\n"
		}
	}
	return componentBlock
}
