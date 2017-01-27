// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import "os"

func filterFiles(files []os.FileInfo, manifestFilename string, result chan string) {
	defer close(result)
	for f := range files {
		if files[f].IsDir() {
			continue
		}
		if files[f].Name() == manifestFilename {
			continue
		}
		result <- files[f].Name()
	}
}
