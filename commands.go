package main

import (
	"github.com/yitsushi/go-commander"
	"log"
	"os"
	"strings"
)

// Parse struct is a representation of the Parse command.
type Parse struct {
}

// Execute is the main function, called on the parse-make command.
func (c *Parse) Execute(opts *commander.CommandHelper) {
	manifest := Manifest{}
	fileName := opts.Arg(0)

	// TODO validate the file to see if we can even open it.

	if isMake(fileName) {
		manifest.parseMake(fileName)
	} else if isLock(fileName) {
		manifest.parseComposer(fileName)
	} else {
		log.Fatal("File type not supported")
		os.Exit(1)
	}

	manifest.compare()
}

func ParseMakefile(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Parse{},
		Help: &commander.CommandDescriptor{
			Name:             "parse",
			ShortDescription: "Parse a Drush makefile",
			LongDescription:  `Parse a Drush makefile (composer.lock file is not supported)`,
			Arguments:        "<filename>",
			Examples: []string{
				"mysite.make",
			},
		},
	}
}

func isLock(filename string) bool {
	return strings.HasSuffix(filename, `.lock`)
}

func isMake(filename string) bool {
	return strings.HasSuffix(filename, `.make`)
}
