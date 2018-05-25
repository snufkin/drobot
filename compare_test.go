package main

import (
	"testing"
)

var testStableVersion = Component{"module", "views", SemVersion{7, 3, 4, ""}}
var testDevVersion = Component{"module", "views", SemVersion{7, 3, -1, "dev"}}
var testGitVersion = Component{"module", "views", SemVersion{7, 3, -1, "git"}}
var testBetaVersionList = []Component{
	{"module", "views", SemVersion{7, 3, 4, "rc1"}},
	{"module", "views", SemVersion{7, 3, 4, "beta5"}},
	{"module", "views", SemVersion{7, 3, 4, "alpha5"}},
}

func TestIsGit(t *testing.T) {
	if ok := testGitVersion.isGit(); !ok {
		t.Error("For", testGitVersion.Version, "expected", true, "got", ok)
	}
}

func TestIsBeta(t *testing.T) {
	for _, v := range testBetaVersionList {
		if ok := v.isBeta(); !ok {
			t.Error("For", v.Version, "expected", true, "got", ok)
		}
	}

	if ok := testGitVersion.isBeta(); ok { // A negative case.
		t.Error("For", testGitVersion.Version, "expected", false, "got", ok)
	}
}

func TestIsDev(t *testing.T) {
	if ok := testDevVersion.isDev(); !ok {
		t.Error("For", testDevVersion.Version, "expected", true, "got", ok)
	}
	if ok := testGitVersion.isStable(); ok {
		t.Error("For", testGitVersion.Version, "expected", false, "got", ok)
	}
}

func TestIsStable(t *testing.T) {
	if ok := testStableVersion.isStable(); !ok {
		t.Error("For", testStableVersion.Version, "expected", true, "got", ok)
	}
	if ok := testGitVersion.isStable(); ok {
		t.Error("For", testGitVersion.Version, "expected", false, "got", ok)
	}
}
