package main

import (
	"fmt"
	"strings"
)

type Component struct {
	// Component type, e.g. module, core etc.
	Type    string
	Name    string
	Version SemVersion
}

type SemVersion struct {
	Major int // Used as the core version for contrib
	Minor int
	Patch int
	Tag   string
}

type Manifest struct {
	Components []Component
}

const CORE = "core"
const MODULE = "module"
const THEME = "theme"

// Convert a semantic version to the d.o format.
func (V SemVersion) printVersion(componentType string, majorVersion int) string {
	if V.Major < 0 || V.Minor < 0 || V.Patch < 0 {
		return fmt.Sprintf("Invalid version")
	}
	switch componentType {
	case CORE:
		if majorVersion == 8 {
			return fmt.Sprintf("%d.%d.%d", V.Major, V.Minor, V.Patch)
		} else {
			return fmt.Sprintf("%d.%d", V.Major, V.Minor)
		}
	case MODULE:
		return fmt.Sprintf("%d.x-%d.%d", V.Major, V.Minor, V.Patch)
	case THEME:
		return fmt.Sprintf("%d.x-%d.%d", V.Major, V.Minor, V.Patch)
	}
	return fmt.Sprintf("%d.x-%d.%d", V.Major, V.Minor, V.Patch)
}

func (C Component) printVersion() string {
	return C.Version.printVersion(C.Type, C.Version.Major)
}

func parseVersion(rawVersion string, majorVersion int) SemVersion {
	version := new(SemVersion)
	if majorVersion == 7 {
		fmt.Sscanf(rawVersion, "%d.%d", &version.Major, &version.Minor)
	}
	return *version
}

func (V *SemVersion) initCore(rawVersion string) {
	parts := strings.Split(rawVersion, ".")

	if len(parts) < 1 || len(parts) > 3 { // Invalid input parses to -1
		V.Major, V.Minor, V.Patch = -1, -1, -1
	} else if len(parts[0]) > 1 {
		V.Major, V.Minor, V.Patch = -1, -1, -1
	} else if len(parts) == 2 {
		fmt.Sscanf(parts[0], "%d", &V.Major)
		fmt.Sscanf(parts[1], "%d", &V.Minor)
	} else if len(parts) == 3 {
		fmt.Sscanf(parts[0], "%d", &V.Major)
		fmt.Sscanf(parts[1], "%d", &V.Minor)
		fmt.Sscanf(parts[2], "%d", &V.Patch)
	}
}

// Parser behaves differently for different core versions.
func (V *SemVersion) initContrib(coreVersion int, rawVersion string) {
	parts := strings.Split(rawVersion, ".")

	if len(parts) < 1 || len(parts) > 3 { // Invalid input parses to -1
		V.Major, V.Minor, V.Patch = -1, -1, -1
		return
	}

	V.Major = coreVersion

	if coreVersion == 7 {
		patch := strings.Split(parts[1], "-")

		// When no patch version is pinned.
		if patch[0] == "x" {
			V.Patch = -1
		} else {
			fmt.Sscanf(patch[0], "%d", &V.Patch)
		}
		fmt.Sscanf(parts[0], "%d", &V.Minor)
		if len(patch) == 2 {
			fmt.Sscanf(patch[1], "%s", &V.Tag)
		}
	} else { // Core: 8, parse the semver from composer.lock (patch is discarded there).
		patch := strings.Split(parts[2], "-")
		fmt.Sscanf(parts[0], "%d", &V.Minor)
		fmt.Sscanf(parts[1], "%d", &V.Patch)

		if len(patch) == 2 {
			fmt.Sscanf(patch[1], "%s", &V.Tag)
		}
	}
}

func (C *Component) init(rawBlock string) {

}

// Collect a manifest list.
func (M *Manifest) append(component Component) {
}
