package main

import (
	"fmt"
	"regexp"
	"strconv"
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
const MANIFEST_MAKE = "make"
const MANIFEST_COMPOSER = "json"

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
	var versionMatches = map[string]string{
		"PARTIAL-STABLE":   `^(\d+)\.(\d+)$`,              // E.g. 3.2
		"PARTIAL-TEST":     `^(\d+)\.(\d+)-(\S+)$`,        // E.g. 4.0-beta7
		"PARTIAL-DEV":      `^(\d+)\.x-dev$`,              // E.g. 5.x-dev
		"DRUPAL-STABLE":    `^\d+\.x-(\d+)\.(\d+)$`,       // E.g. 7.x-4.3
		"DRUPAL-TEST":      `^\d+\.x-(\d+)\.(\d+)-(\S+)$`, // E.g. 7.x-4.3-beta1
		"DRUPAL-DRUSH-DEV": `^\d+\.x-(\d+)\.x$`,           // E.g. 7.x-2.x
		"DRUPAL-DEV":       `^\d+\.x-(\d+)\.x-dev$`,       // E.g. 7.x-2.x-dev
		"COMPOSER-STABLE":  `^(\d+)\.(\d+)\.\d+$`,         // E.g. 1.1.0
		"COMPOSER-TEST":    `^(\d+)\.(\d+)\.\d+-(\S+)$`,   // E.g. 1.1.0-beta1
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
			case "PARTIAL-DEV", "DRUPAL-DRUSH-DEV", "DRUPAL-DEV":
				V.Minor, _ = strconv.Atoi(matches[1])
				V.Patch, V.Tag = -1, "dev"
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
func (C *Component) init(coreVersion int, rawVersion string, componentType string) {
	switch componentType {
	case CORE:
		C.Version.initCore(rawVersion)
	case MODULE:
	case THEME:
		C.Version.initContrib(coreVersion, rawVersion)
	}
}

// Initialise a manifest list out of the raw txt from the yaml file.
func (M *Manifest) initFromMake(rawBlock string) {
	//1. Identify the names of the projects
	//2. Extract all variations of keys (assumption and may break)
	componentList := componentList(rawBlock)

	for _, component := range componentList {
		block := findBlock(component)
		println(block)
	}
	// Find the string block that contains each components and initialise the component.

}

// Collect a manifest list.
func (M *Manifest) append(component Component) {
}
