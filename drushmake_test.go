package main

import (
	"testing"
)

type versionTest struct {
	Expected string
	Actual   string
	Error    error
}

var testCoreVersions = []struct {
	in      string
	out     SemVersion
	success bool
}{
	// Core version denominations.
	{"6.13", SemVersion{6, 13, 0, ""}, true},
	{"6.13", SemVersion{5, 13, 0, ""}, false},
	{"7.0", SemVersion{7, 0, 0, ""}, true},
	{"8.3.4", SemVersion{8, 3, 4, ""}, true},
	{"8.0.0", SemVersion{8, 0, 0, ""}, true},
	{"88.00.00", SemVersion{-1, -1, -1, ""}, true},
	{"88.00.00", SemVersion{88, 0, 0, ""}, false},
}

var testContribVersions = []struct {
	core    int // Since D7 uses .make, D8 uses composer parser needs to know.
	in      string
	out     SemVersion
	success bool
}{
	{7, "2.0beta1", SemVersion{7, 2, 0, ""}, false}, // No dash.
	{7, "1.0-rc1", SemVersion{7, 1, 0, "rc1"}, true},
	{7, "2.x-dev", SemVersion{7, 2, -1, "dev"}, true},
	{7, "2.x-dev", SemVersion{7, 2, -1, ""}, false},
	{7, "2.x-dev", SemVersion{7, 2, 0, ""}, false},
	{8, "1.0.0-rc2", SemVersion{8, 1, 0, "rc2"}, true},
	{8, "1.1.0", SemVersion{8, 1, 1, ""}, true},
}

// Project was cloned from git and is not pinned to a hash.
var componentBlockGit = `projects[ns_core][type] = module
projects[ns_core][download][type] = git
projects[ns_core][download][branch] = 7.x-2.x`

// Project was cloned from git and it is pinned to a hash.
var componentBlockGitHash = `projects[draggableviews][type] = module
projects[draggableviews][download][type] = git
projects[draggableviews][download][revision] = 9677bc18b7255e13c33ac3cca48732b855c6817d
projects[draggableviews][download][branch] = 7.x-2.x`

// Project is pinned to a stable version on a single line. Test assumes to know core version.
var componentBlockVersionOneLine = `
projects[views] = 3.1`

// Project is pinned to a stable version on multiple lines.
var componentBlockVersionMultiLine = `
projects[nodequeue][subdir] = contrib
projects[nodequeue][version] = 2.0-alpha1
projects[nodequeue][patch][] = "http://drupal.org/files/issues/1023606-qid-to-name-6.patch"
projects[nodequeue][patch][] = "http://drupal.org/files/issues/nodequeue_d7_autocomplete-872444-6.patch"`

var testComponentBlocks = []struct {
	core    int
	in      string
	out     Component
	success bool
}{
	{
		core: 7,
		in:   componentBlockGit,
		out: Component{
			Type:    "module",
			Name:    "ns_core",
			Version: SemVersion{7, 2, -1, "git"},
		},
		success: true,
	},
	{
		core: 7,
		in:   componentBlockGitHash,
		out: Component{
			Type:    "module",
			Name:    "draggableviews",
			Version: SemVersion{7, 2, -1, "git:9677bc18b7255e13c33ac3cca48732b855c6817d"},
		},
		success: true,
	},
	{
		core: 7,
		in:   componentBlockVersionOneLine,
		out: Component{
			Type:    "module",
			Name:    "views",
			Version: SemVersion{7, 3, 1, ""},
		},
		success: true,
	},
	{
		core: 7,
		in:   componentBlockVersionMultiLine,
		out: Component{
			Type:    "module",
			Name:    "nodequeue",
			Version: SemVersion{7, 2, 0, "alpha1"},
		},
		success: true,
	},
}

// Test if component blocks are correctly parsed and populated.
func TestComponentBlockParser(t *testing.T) {
	for _, testBlock := range testComponentBlocks {
		b := Component{}
		b.init(testBlock.in)
		if success := (b == testBlock.out); b != testBlock.out && success != testBlock.success {
			t.Error("For", testBlock.in, "expected", testBlock.out, "got", b)
		}

	}
}

// Test if versions are correctly translated into the semantic structure.
func TestCoreVersionParser(t *testing.T) {
	for _, testV := range testCoreVersions {
		v := SemVersion{}
		v.initCore(testV.in)
		if success := (v == testV.out); v != testV.out && success != testV.success {
			t.Error(
				"For", testV.in,
				"expected", testV.out,
				"got", v,
			)
		}
	}
}

// Test if versions are correctly translated into the semantic structure.
func TestContribVersionParser(t *testing.T) {
	for _, testV := range testContribVersions {
		v := SemVersion{}
		v.initContrib(testV.core, testV.in)
		if success := (v == testV.out); v != testV.out && success != testV.success {
			t.Error(
				"For", testV.in,
				"expected", testV.out,
				"got", v,
			)
		}
	}
}

func TestParseMakefile(t *testing.T) {
	parseMakefile("test/test.make")
}
