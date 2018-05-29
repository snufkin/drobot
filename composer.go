package main

import (
	"fmt"
	clp "github.com/snufkin/go-composerlockparser"
	"strings"
)

func (M *Manifest) parseComposer(path string) {
	lock := clp.ComposerInfo{PathToLockFile: path}
	lock.Parse()
	pList := lock.GetPackageListByPrefix("drupal")

	for _, p := range pList {
		c := Component{}
		c.convertToComponent(p)
		M.Components = append(M.Components, c)
	}
}

func (c *Component) convertToComponent(p clp.Package) {
	// Clean up the project name.
	c.Name = strings.Replace(p.Name, "drupal/", "", 1)

	if p.Type == "drupal-core" {
		c.Type = "core"
		fmt.Sscanf(p.Version, "%d.%d.%d", &c.Version.Major, &c.Version.Minor, &c.Version.Patch)
	} else {
		c.Type = "module"
		c.Version.Major = 8
		fmt.Sscanf(p.Version, "%d.%d.%d", &c.Version.Minor, &c.Version.Patch)
	}
}
