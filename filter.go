// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

package main

// filterFiles outputs only files that are not the manifest
func filterFiles(done chan struct{}, files chan *pathFileInfo, manifestFilename string, out chan *pathFileInfo) {
	defer close(out)
	for file := range files {
		select {
		case <-done:
			return
		default:
			if filterFile(file, manifestFilename) {
				out <- file
			}
		}

	}
	return
}

// filterFile returns true if file should be hashed
func filterFile(file *pathFileInfo, manifestFileName string) bool {
	if file.IsDir() {
		return false
	}
	if file.Name() == manifestFileName {
		return false
	}
	if file.name == manifestFileName {
		return false
	}
	return true

}
