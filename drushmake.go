package main

import (
	dmp "github.com/snufkin/go-drushmakeparser"
)

func (M *Manifest) parseMake(path string) {
	makeInfo := dmp.DrushMakeInfo{}
	makeInfo.Parse(path)

	for _, p := range makeInfo.Packages {
		c := Component{}
		if p.Version == "" && p.Download.Branch != "" {
			c.init(7, p.Download.Branch, p.Name, p.Type)
		} else {
			c.init(7, p.Version, p.Name, p.Type)
		}
		M.Components = append(M.Components, c)
	}
}
