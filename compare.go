package main

import (
	"fmt"
)

// Possible outcome of a version comparison.
// When a stable release is available, but we using an earlier, or dev version.
const STABLE_AVAILABLE = 0

// When both modules are stable and the current one is outdated.
const OUTDATED = 1

// When the module is up to date.
const UPDATED = 2

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
	case UPDATED:
		fmt.Printf("[OK]\t%s", C.Name)
	case OUTDATED:
		fmt.Printf("[UPDATE AVAILABLE]\t%s", C.Name)
	case STABLE_AVAILABLE:
		fmt.Printf("[STABLE AVAILABLE]\t%s", C.Name)
	}
	fmt.Printf("\tActual: %v\tCurrent: %v\n", C.Version, release.Releases[0])
}

func (C Component) checkUpdateStatus(r Release) int {
	if r.Patch > C.Version.Patch {
		return OUTDATED
		// fmt.Printf("Component %s is outdated %v => %v\n", C.Name, C.Version, latestRelease)
	} else if r.Minor == C.Version.Minor && r.Patch == C.Version.Patch {
		return UPDATED
	} else if C.Version.Patch == -1 && C.Version.Minor == r.Minor {
		return STABLE_AVAILABLE
	}
	return -1
}

// Compare two versions and return a status evaluation.
func (C Component) status() {

}
