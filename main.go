// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

// Helps create a new `manifest.json` file in a folder, and verify an existing manifest.
// Each file is checked for `SHA1` and `MD5` hashes.
//
// This is helpful to see if the contents of a file have changed since the last time the tool was run.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const VERIFY_MANIFEST_VERSION = "v0.2"
const VERIFY_MANIFEST_WEBSITE = "https://github.com/robert-wallis/VerifyManifest"

type commandFlag struct {
	RootDir          string
	ManifestFilename string
}

var g_flags = commandFlag{}

func init() {
	flag.StringVar(&g_flags.RootDir, "root", ".", "Root folder to calculate Sum.")
	flag.StringVar(&g_flags.ManifestFilename, "manifest", "manifest.json", "Manifest file name.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\nVersion %s\n%s\n\n", os.Args[0], VERIFY_MANIFEST_VERSION, VERIFY_MANIFEST_WEBSITE)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	infoLog := log.New(os.Stdout, "", 0)
	errorLog := log.New(os.Stderr, "", 0)
	hashFolder(g_flags.RootDir, g_flags.ManifestFilename, infoLog, errorLog)
}
