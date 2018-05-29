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
	Tag   string `xml:"version_extra"`
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

func (r Release) String() string {
	if r.Tag == "" {
		return fmt.Sprintf("%d.%d.%d", r.Major, r.Minor, r.Patch)
	}
	return fmt.Sprintf("%d.%d.%d-%s", r.Major, r.Minor, r.Patch, r.Tag)
}

func (C *Component) fetchReleases() Result {
	url := fmt.Sprintf("%s/%s/%d.x", RELEASE_URL, C.Name, C.Version.Major)

	r := Result{Name: "", Type: ""}

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
		fmt.Printf("marshal error: %v", err)
	}

	// If the release was NOT core, then rearrange the structure.
	// I really did not want to write a custom unmarhsaller.
	if C.Name != "drupal" {
		for k, rls := range r.Releases {
			r.Releases[k].Minor, r.Releases[k].Major = rls.Major, C.Version.Major
		}
	}

	return r
}
