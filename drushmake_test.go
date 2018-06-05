package main

import (
	"testing"
)


var testCoreVersions = []struct {
	in      string
	out     SemVersion
	success bool
}{
	// Core version denominations.
	{"6.13", SemVersion{6, 0, 13, ""}, true},
	{"6.13", SemVersion{5, 0, 13, ""}, false},
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
	{7, "3.4", SemVersion{7, 3, 4, ""}, true},
	{7, "2.0beta1", SemVersion{7, 2, 0, ""}, false}, // No dash, no success.
	{7, "1.0-rc1", SemVersion{7, 1, 0, "rc1"}, true},
	{7, "2.x-dev", SemVersion{7, 2, -1, "dev"}, true},
	{7, "2.x-dev", SemVersion{7, 2, -1, ""}, false},
	{7, "2.x-dev", SemVersion{7, 2, 0, ""}, false},
	{8, "1.0.0-rc2", SemVersion{8, 1, 0, "rc2"}, true},
	{8, "1.1.0", SemVersion{8, 1, 1, ""}, true},
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
		v.initContribVersion(testV.core, testV.in)
		if success := (v == testV.out); v != testV.out && success != testV.success {
			t.Error(
				"For", testV.in,
				"expected", testV.out,
				"got", v,
			)
		}
	}
}
