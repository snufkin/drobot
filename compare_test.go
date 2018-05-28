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

var statusList = []struct {
	release   Release
	component Component
	out       int
	success   bool
}{
	{
		release:   Release{Major: 7, Minor: 3, Patch: 4, Tag: ""},
		component: Component{"module", "views", SemVersion{7, 3, 4, ""}},
		out:       OK,
		success:   true,
	},
	{
		release:   Release{Major: 7, Minor: 3, Patch: 4, Tag: ""},
		component: Component{"module", "views", SemVersion{7, 3, 3, ""}},
		out:       UPDATE_AVAILABLE,
		success:   true,
	},
	{
		release:   Release{Major: 7, Minor: 3, Patch: 4, Tag: ""},
		component: Component{"module", "views", SemVersion{7, 2, 4, ""}},
		out:       STABLE_MAJOR_AVAILABLE,
		success:   true,
	},
	{
		release:   Release{Major: 7, Minor: 3, Patch: 1, Tag: ""},
		component: Component{"module", "views", SemVersion{7, 3, -1, "dev"}},
		out:       STABLE_AVAILABLE,
		success:   true,
	},
	{
		release:   Release{Major: 7, Minor: 3, Patch: 1, Tag: ""},
		component: Component{"module", "views", SemVersion{7, 3, -1, "git"}},
		out:       STABLE_AVAILABLE,
		success:   true,
	},
}

func TestCheckUpdateStatus(t *testing.T) {
	for _, testCase := range statusList {
		if status := testCase.component.checkUpdateStatus(testCase.release); (status == testCase.out) != testCase.success {
			t.Error("For", testCase.component, "vs", testCase.release, "expected", messages[testCase.out].short, "got", messages[status].short)
		}
	}
}
