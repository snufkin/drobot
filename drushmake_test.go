package main

import (
	"testing"
)

type versionTest struct {
	Expected string
	Actual   string
	Error    error
}

var testVersions = []struct {
	in      string
	out     SemVersion
	success bool
}{
	{"6.13", SemVersion{6, 13, 0}, true},
	{"6.13", SemVersion{5, 13, 0}, false},
	{"7.0", SemVersion{7, 0, 0}, true},
	{"8.3.4", SemVersion{8, 3, 4}, true},
	{"8.0.0", SemVersion{8, 0, 0}, true},
	{"88.00.00", SemVersion{-1, -1, -1}, true},
	{"88.00.00", SemVersion{88, 0, 0}, false},
}

func TestCoreVersionParser(t *testing.T) {
	for _, testV := range testVersions {
		v := SemVersion{}
		v.init(testV.in)
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
