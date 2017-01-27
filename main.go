// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

// Helps create a new `manifest.json` file in a folder, and verify an existing manifest.
// Each file is checked for `SHA1` and `MD5` hashes.
//
// This is helpful to see if the contents of a file have changed since the last time the tool was run.
package main

import (
	"flag"
	"log"
	"os"
)

var rootDir = flag.String("root", ".", "Root folder to calculate Sum.")
var manifestFilename = flag.String("manifest", "manifest.json", "Manifest file name.")

func main() {
	flag.Parse()
	infoLog := log.New(os.Stdout, "", 0)
	errorLog := log.New(os.Stderr, "", 0)
	hashFolder(*rootDir, *manifestFilename, infoLog, errorLog)
}
