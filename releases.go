package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Term struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

// Information about a single release.
type Release struct {
	Major int    `xml:"version_major"`
	Minor int    `xml:"version_minor"`
	Patch int    `xml:"version_patch"`
	Terms []Term `xml:"terms>term"`
}

// The map of the release history response.
type Result struct {
	XMLName          xml.Name  `xml:"project"`
	Name             string    `xml:"short_name"`
	RecommendedMajor int       `xml:"recommended_major"`
	Type             string    `xml:"type"`
	Releases         []Release `xml:"releases>release"`
}

// Example: https://updates.drupal.org/release-history/drupal/7.x
const RELEASE_URL = `https://updates.drupal.org/release-history`

func fetchRelease(cName string, coreVersion int) {
	// url := fmt.Sprintf("%s/%s/%d.x", RELEASE_URL, cName, coreVersion)

	r := Result{Name: "", Type: ""}

	data, err := ioutil.ReadFile("test/drupal-7.x.xml")

	if err != nil {
		fmt.Printf("XML File read error: %v", err)
		return
	}

	err = xml.Unmarshal([]byte(data), &r)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	if cName == "drupal" {
		fmt.Printf("Name: %#v\n", r.Name)
		fmt.Printf("Type: %#v\n", r.Type)
		fmt.Printf("Recommended major: %#v\n", r.RecommendedMajor)
		fmt.Printf("Release history: %#v\n", r.Releases)
	}
}
