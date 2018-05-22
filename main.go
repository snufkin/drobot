package main

import (
	"fmt"
)

func main() {
	manifest := Manifest{}
	manifest.parseMakefile("test/test.make")
	for _, c := range manifest.Components {
		fmt.Printf("%s => %v\n", c.Name, c.Version)
	}
}
