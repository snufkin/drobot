package main

import (
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
	fmt.Println(componentList)

	// Grab core information
	rawCoreVersion := cfg.Section("").Key("core")
	// core := Component{CORE, "drupal", parseVersion(rawCoreVersion.Value(), 7)}

	version := SemVersion{}
	version.init(rawCoreVersion.Value())
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
