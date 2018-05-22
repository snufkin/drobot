package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func fetchRelease(cName string, coreVersion int) Result {
	url := fmt.Sprintf("%s/%s/%d.x", RELEASE_URL, cName, coreVersion)

	r := Result{Name: "", Type: ""}

	// data, err := ioutil.ReadFile("test/drupal-7.x.xml")

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal([]byte(data), &r)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	// If the release was NOT core, then rearrange the structure.
	// I really did not want to write a custom unmarhsaller.
	if cName != "drupal" {
		for k, c := range r.Releases {
			r.Releases[k].Minor = c.Major
			r.Releases[k].Major = coreVersion
		}
	}

	return r

	// if cName == "drupal" {
	// 	fmt.Printf("Name: %#v\n", r.Name)
	// 	fmt.Printf("Type: %#v\n", r.Type)
	// 	fmt.Printf("Recommended major: %#v\n", r.RecommendedMajor)
	// 	fmt.Printf("Release history: %#v\n", r.Releases)
	// }
}

// Check the update status for a given manifest element.
func (C Component) checkUpdate() {
	release := fetchRelease(C.Name, C.Version.Major)
	latestRelease := release.Releases[0]
	if latestRelease.Patch > C.Version.Patch {
		fmt.Printf("Component %s is outdated %v => %v\n", C.Name, C.Version, latestRelease)
	}
}

// Compare two versions and return a status evaluation.
func (C Component) status() {

}
