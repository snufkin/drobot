package main

// Check each of the manifest elements against the release data.
func (M Manifest) compare() {
	for _, c := range M.Components {
		c.checkUpdate()

	}

}
