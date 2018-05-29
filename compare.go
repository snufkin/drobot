package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Possible outcome of a version comparison.
const (
	UNKNOWN                = -1 // Case not handled.
	OK                     = 2  // When the module is up to date.
	STABLE_AVAILABLE       = 0  // Stable release is available, when using non-stable.
	UPDATE_AVAILABLE       = 1  // Both modules are stable and the current one is outdated.
	STABLE_MAJOR_AVAILABLE = 6  // Major version upgrade is available.
	UNSUPPORTED_MAJOR      = 7  // The major version of the used version is no longer supported.
	BETA_AVAILABLE         = 8  // Beta release available, when using dev.
)

type message struct {
	short string
	long  string
}

var messages = map[int]message{
	UNKNOWN:                {short: "UNKOWN", long: "Version difference not implemented"},
	OK:                     {short: "OK", long: "No update required."},
	STABLE_AVAILABLE:       {short: "UPDATE AVAILABLE", long: "Stable release is available, while using non-pinned, update advised."},
	UPDATE_AVAILABLE:       {short: "UPDATE AVAILABLE", long: "Update is available."},
	STABLE_MAJOR_AVAILABLE: {short: "MAJOR UPDATE AVAILABLE", long: "Major update is available"},
	UNSUPPORTED_MAJOR:      {short: "NOT SUPPORTED", long: "Used major version is no longer supported."},
	BETA_AVAILABLE:         {short: "BETA UPDATE AVAILABLE", long: "Beta release available, update advised."},
}

// Check each of the manifest elements against the release data.
func (M Manifest) compare() {
	fmt.Printf("Status\tDescription\tModule\tCurrent\tAvailable\n")
	for _, c := range M.Components {
		c.checkUpdate()
	}
}

// Check the update status for a given manifest element. TODO don't print here cmon.
func (C Component) checkUpdate() {
	releases := C.fetchReleases()
	status := C.checkUpdateStatus(releases.Releases[0])

	fmt.Printf("[%s]\t%s\t%s\t%s\t%s\n", messages[status].short, messages[status].long, C.Name, C.Version, releases.Releases[0])
}

// Compare two versions and assign a status to the component.
func (C Component) checkUpdateStatus(r Release) int {
	// No pinned version is available.
	if r.Tag == "dev" {
		if C.Version.Tag != "" {

		}
	}

	// When stable is available but is not used.
	if !C.isStable() && r.Tag == "" && &r.Minor != nil {
		return STABLE_AVAILABLE
	} else if C.isStable() && r.Minor > C.Version.Minor { // Minor update available.
		return STABLE_MAJOR_AVAILABLE
	} else if C.isStable() && r.Minor == C.Version.Minor && r.Patch > C.Version.Patch {
		return UPDATE_AVAILABLE
	} else if r.Minor == C.Version.Minor && r.Patch == C.Version.Patch { // Same version.
		return OK
	} else if C.isDev() && rIsBeta(r.Tag) { // Beta available.
		return BETA_AVAILABLE
	} else if C.isBeta() && rIsBeta(r.Tag) { // Both releases are betas.

		re := regexp.MustCompile("^(rc|beta|alpha)([0-9]+)$")

		currentTagMatches, releaseTagMatches := re.FindStringSubmatch(C.Version.Tag), re.FindStringSubmatch(r.Tag)
		if len(currentTagMatches) == 3 && len(releaseTagMatches) == 3 {
			currentTag, currentVersion := currentTagMatches[1], currentTagMatches[2]
			releaseTag, releaseVersion := releaseTagMatches[1], releaseTagMatches[2]

			if currentTag == releaseTag && releaseVersion > currentVersion {
				return BETA_AVAILABLE
			} else if currentTag == "alpha" && (releaseTag == "rc" || releaseTag == "beta") {
				return BETA_AVAILABLE
			} else if currentTag == "beta" && releaseTag == "rc" {
				return BETA_AVAILABLE
			}
		}
	} else if C.isGit() && r.Tag != "dev" {
		return BETA_AVAILABLE
	}

	return UNKNOWN
}

// Shortcut to a beta release.
func rIsBeta(tag string) bool {
	return tag != "" && tag != "dev"
}

func (C Component) isGit() bool {
	return C.Version.Tag == "git"
}

// Non-stable, but fixed release.
func (C Component) isBeta() bool {
	fixedTags := []string{"rc", "beta", "alpha"}
	for _, s := range fixedTags {
		if strings.HasPrefix(C.Version.Tag, s) {
			return true
		}
	}
	return false
}

func (C Component) isDev() bool {
	return C.Version.Tag == "dev"
}

func (C Component) isStable() bool {
	return C.Version.Tag == ""
}

// Compare two versions and return a status evaluation.
func (C Component) status() {

}
