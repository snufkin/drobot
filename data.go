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
	Major int64 // Used as the core version for contrib
	Minor int64
	Patch int64
}

type Manifest struct {
	Components []Component
}

const CORE = "core"
const MODULE = "module"
const THEME = "theme"

// Convert a semantic version to the d.o format.
func (V SemVersion) printVersion(componentType string, majorVersion int64) string {
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

func parseVersion(rawVersion string, majorVersion int64) SemVersion {
	version := new(SemVersion)
	if majorVersion == 7 {
		fmt.Sscanf(rawVersion, "%d.%d", &version.Major, &version.Minor)
	}
	return *version
}

func (V *SemVersion) init(rawVersion string) {
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

// Collect a manifest list.
func (M *Manifest) append(component Component) {
}
