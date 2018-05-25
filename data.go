package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Component type, e.g. module, core etc.
type Component struct {
	Type    string
	Name    string
	Version SemVersion
}

// Storage for the converted semantic-ish version of the component.
type SemVersion struct {
	Major int // Note: core version for contrib.
	Minor int
	Patch int    // Note: D7 uses Patch (e.g. 7.44 is 7.0.44)
	Tag   string // Additional information such as rc1, git, dev.
}

// Manifest holds the complete list of components as parsed from the file.
type Manifest struct {
	Components []Component
}

const (
	CORE              = "core"   // Represents Drupal core itself.
	MODULE            = "module" // Plugin type is module.
	THEME             = "theme"  // Plugin type is theme.
	MANIFEST_MAKE     = "make"   // The manifest file is in drush make format.
	MANIFEST_COMPOSER = "json"   // The manifest file is a composer.lock file.
)

// Convert a semantic version to the d.o format.
func (V SemVersion) printVersion(componentType string, majorVersion int) string {
	if V.Major < 0 && V.Minor < 0 && V.Patch < 0 {
		return fmt.Sprintf("Invalid version")
	}
	switch componentType {
	case CORE:
		if majorVersion == 8 {
			return fmt.Sprintf("%d.%d.%d", V.Major, V.Minor, V.Patch)
		}
		return fmt.Sprintf("%d.%d", V.Major, V.Patch)
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

func (V SemVersion) String() string {
	if V.Tag == "" {
		return fmt.Sprintf("%d.%d.%d", V.Major, V.Minor, V.Patch)
	}
	return fmt.Sprintf("%d.%d.%d-%s", V.Major, V.Minor, V.Patch, V.Tag)
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
		fmt.Sscanf(parts[1], "%d", &V.Patch) // 7.44 resolves to 7.0.44
	} else if len(parts) == 3 {
		fmt.Sscanf(parts[0], "%d", &V.Major)
		fmt.Sscanf(parts[1], "%d", &V.Minor)
		fmt.Sscanf(parts[2], "%d", &V.Patch)
	}
}

// Parser behaves differently for different core versions.
func (V *SemVersion) initContribVersion(coreVersion int, rawVersion string) {
	var versionMatches = map[string]string{
		"PARTIAL-STABLE":  `^(\d+)\.(\d+)$`,              // E.g. 3.2
		"PARTIAL-TEST":    `^(\d+)\.(\d+)-(\w+)$`,        // E.g. 4.0-beta7
		"PARTIAL-DEV":     `^(\d+)\.x-(\w+)$`,            // E.g. 5.x-dev
		"DRUPAL-STABLE":   `^\d+\.x-(\d+)\.(\d+)$`,       // E.g. 7.x-4.3
		"DRUPAL-TEST":     `^\d+\.x-(\d+)\.(\d+)-(\w+)$`, // E.g. 7.x-4.3-beta1
		"DRUPAL-DRUSH":    `^\d+\.x-(\d+)\.x$`,           // E.g. 7.x-2.x
		"DRUPAL-DEV":      `^\d+\.x-(\d+)\.x-(dev|git)$`, // E.g. 7.x-2.x-dev/git (prepared)
		"GIT-HASH":        `^\d+\.x-(\d+)\.x-(git:\w+)$`, // E.g. 7.x-2.x-git:<long hash> (prepared)
		"COMPOSER-STABLE": `^(\d+)\.(\d+)\.\d+$`,         // E.g. 1.1.0
		"COMPOSER-TEST":   `^(\d+)\.(\d+)\.\d+-(\w+)$`,   // E.g. 1.1.0-beta1
	}

	V.Major, V.Tag = coreVersion, ""

	foundMatch := false
	for vType, expression := range versionMatches {
		re := regexp.MustCompile(expression)
		if isMatch := re.MatchString(rawVersion); isMatch {
			matches := re.FindStringSubmatch(rawVersion)

			switch vType {
			case "PARTIAL-STABLE", "DRUPAL-STABLE", "COMPOSER-STABLE":
				V.Minor, _ = strconv.Atoi(matches[1])
				V.Patch, _ = strconv.Atoi(matches[2])
				foundMatch = true
			case "PARTIAL-TEST", "DRUPAL-TEST", "COMPOSER-TEST":
				V.Minor, _ = strconv.Atoi(matches[1])
				V.Patch, _ = strconv.Atoi(matches[2])
				V.Tag = matches[3]
				foundMatch = true
			case "DRUPAL-DRUSH":
				V.Minor, _ = strconv.Atoi(matches[1])
				V.Patch = -1
				foundMatch = true
			case "PARTIAL-DEV", "DRUPAL-DEV", "GIT-HASH":
				V.Minor, _ = strconv.Atoi(matches[1])
				V.Patch, V.Tag = -1, matches[2]
				foundMatch = true
			}
		}
	}

	if !foundMatch {
		V.Major, V.Minor, V.Patch = -1, -1, -1
	}
	return
}

// Prepare a component with various conditional initialisers.
func (C *Component) init(coreVersion int, rawVersion string, componentName string, componentType string) {
	switch componentType {
	case CORE:
		C.Version.initCore(rawVersion)
		C.Name = componentName
		C.Type = componentType
	case MODULE:
		fallthrough
	case THEME:
		C.Version.initContribVersion(coreVersion, rawVersion)
		C.Name = componentName
		C.Type = componentType
	}
}
