package main

import (
	"github.com/yitsushi/go-commander"
)

func main() {
	registry := commander.NewCommandRegistry()
	registry.Register(ParseMakefile)
	registry.Execute()

}

// Parse struct is a representation of the Parse command.
type Parse struct {
}

// Execute is the main function, called on the parse-make command.
func (c *Parse) Execute(opts *commander.CommandHelper) {
	manifest := Manifest{}
	fileName := opts.Arg(0)
	// Todo validate the file.
	manifest.parseMakefile(fileName)
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
