package main

import (
	"fmt"
	"github.com/go-ini/ini"
)

// Parse a makefile on filePath and return the list of modules, themes and core.
func parseMakefile(filePath string) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
	}

	// keys := cfg.Section("").Keys()
	// names := cfg.Section("").KeyStrings()
	// hash := cfg.Section("").KeysHash()

	// for key, val := range hash {
	// fmt.Printf("%v => %v\n", key, val)
	// }

	// Grab core information
	rawCoreVersion := cfg.Section("").Key("core")
	// core := Component{CORE, "drupal", parseVersion(rawCoreVersion.Value(), 7)}

	version := SemVersion{}
	version.init(rawCoreVersion.Value())
	// core := Component{CORE, "drupal", version}

	// fmt.Printf("%v\n", core.printVersion())
}
