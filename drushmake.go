package main

import (
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
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
	C.Name = keyMapper(block)

	scanner := bufio.NewScanner(strings.NewReader(block))
	for scanner.Scan() {
		// Split the string along the = symbol.
		lineParts := strings.Split(scanner.Text(), " = ")

		// Multi-line version definition
		if key := strings.Replace(lineParts[0], "projects["+C.Name+"]", "", 1); key != "" {
		} else {
			C.init(coreVersion, lineParts[1], MODULE)
		}
	}

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
