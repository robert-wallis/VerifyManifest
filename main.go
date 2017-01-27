// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

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
