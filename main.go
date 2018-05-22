package main

func main() {
	manifest := Manifest{}
	manifest.parseMakefile("test/test.make")
	manifest.compare()

	// for _, c := range manifest.Components {
	// fetchRelease(c.Name, c.Version.Major)
	// }
}
