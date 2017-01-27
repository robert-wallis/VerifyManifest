// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

package main

import (
	"flag"
)

var rootDir = flag.String("root", ".", "Root folder to calculate Sum.")
var manifestFilename = flag.String("manifest", "manifest.json", "Manifest file name.")

func main() {
	flag.Parse()
	hashFolder(*rootDir, *manifestFilename)
}
