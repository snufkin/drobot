package main

import (
	"github.com/yitsushi/go-commander"
)

func main() {
	registry := commander.NewCommandRegistry()
	registry.Register(ParseMakefile)
	registry.Execute()
}
