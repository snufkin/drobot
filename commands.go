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

func DrobotFileValidator(c *commander.CommandHelper) {
	if c.Arg(0) == "" {
		panic("File not specified")
	}

	info, err := os.Stat(c.Arg(0))
	if err != nil {
		panic("File not found")
	}

	if !info.Mode().IsRegular() {
		panic("Not a regular file, can not process.")
	}
}

// Command definition to initiate the parsing of a manifest file.
func ParseMakefile(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Parse{},
		Validator: DrobotFileValidator,
		Help: &commander.CommandDescriptor{
			Name:             "parse",
			ShortDescription: `Parse a site manifest file and report outdated components.`,
			LongDescription:  `Process a Drush makefile, or composer.lock file and check each component against releases from drupal.org.`,
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
