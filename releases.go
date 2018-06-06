package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Term maps to the taxonomy terms on releases.
type Term struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

// Release holds information about a single release of a package.
type Release struct {
	Major int    `xml:"version_major"`
	Minor int    `xml:"version_minor"`
	Patch int    `xml:"version_patch"`
	Tag   string `xml:"version_extra"`
	Terms []Term `xml:"terms>term"`
}

// Result stores the release history response from drupal.org
type Result struct {
	XMLName          xml.Name  `xml:"project"`
	Name             string    `xml:"short_name"`
	RecommendedMajor int       `xml:"recommended_major"`
	Type             string    `xml:"type"`
	Releases         []Release `xml:"releases>release"`
}

// RELEASE_URL is the base URL where we get the release history from.
// Example pattern: https://updates.drupal.org/release-history/drupal/7.x
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
	if err != nil && resp.ContentLength == 127 {
		fmt.Println("Project not found.")
		r.Releases = []Release{}
	} else if err != nil {
		fmt.Printf("%v Marshal error on url %s: %v\n", C.Name, url, err)
	}

	// If the release was NOT core, then rearrange the structure.
	// I really did not want to write a custom unmarhsaller.
	if C.Name != "drupal" && len(r.Releases) > 0 {
		for k, rls := range r.Releases {
			r.Releases[k].Minor, r.Releases[k].Major = rls.Major, C.Version.Major
		}
	}

	return r
}
