package main

import (
	"fmt"
	"strings"
)

// Possible outcome of a version comparison.
const (
	DEV_STABLE_AVAILABLE   = 0 // Stable release is available (when using a dev)
	STABLE_STATUS_OUTDATED = 1 // Both modules are stable and the current one is outdated.
	STABLE_STATUS_UPDATED  = 2 // When the module is up to date.
)

// Check each of the manifest elements against the release data.
func (M Manifest) compare() {
	for _, c := range M.Components {
		c.checkUpdate()
	}
}

// Check the update status for a given manifest element.
func (C Component) checkUpdate() {
	release := fetchRelease(C.Name, C.Version.Major)
	status := C.checkUpdateStatus(release.Releases[0])
	switch status {
	case STABLE_STATUS_UPDATED:
		fmt.Printf("[OK]\t%s", C.Name)
	case STABLE_STATUS_OUTDATED:
		fmt.Printf("[UPDATE AVAILABLE]\t%s", C.Name)
	case DEV_STABLE_AVAILABLE:
		fmt.Printf("[STABLE AVAILABLE]\t%s", C.Name)
	}
	fmt.Printf("\tActual: %v\tCurrent: %v\n", C.Version, release.Releases[0])
}

func (C Component) checkUpdateStatus(r Release) int {
	if r.Patch > C.Version.Patch {
		return STABLE_STATUS_OUTDATED
		// fmt.Printf("Component %s is outdated %v => %v\n", C.Name, C.Version, latestRelease)
	} else if r.Minor == C.Version.Minor && r.Patch == C.Version.Patch {
		return STABLE_STATUS_UPDATED
	} else if C.Version.Patch == -1 && C.Version.Minor == r.Minor {
		return DEV_STABLE_AVAILABLE
	}
	return -1
}

func (C Component) isGit() bool {
	return C.Version.Tag == "git"
}

// Non-stable, but fixed release.
func (C Component) isBeta() bool {
	fixedTags := []string{"rc", "beta", "alpha"}
	for _, s := range fixedTags {
		if strings.Index(C.Version.Tag, s) == 0 {
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
