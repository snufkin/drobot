package main

import (
	"fmt"
	clp "github.com/snufkin/go-composerlockparser"
	"regexp"
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

	// Check if we have a tag.
	re := regexp.MustCompile(`.*-(rc|beta|alpha)([0-9]+)$`)

	tagMatch := re.FindStringSubmatch(p.Version)

	if len(tagMatch) == 3 {
		c.Version.Tag = fmt.Sprintf("%s-%s", tagMatch[1], tagMatch[2])
	}

	if p.Type == "drupal-core" {
		c.Type = "core"
		fmt.Sscanf(p.Version, "%d.%d.%d", &c.Version.Major, &c.Version.Minor, &c.Version.Patch)
	} else {
		c.Type = "module"
		c.Version.Major = 8
		fmt.Sscanf(p.Version, "%d.%d.%d", &c.Version.Minor, &c.Version.Patch)
	}

}
