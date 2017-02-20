// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

// Command VerifyManifest helps create a new `manifest.json` file in a folder, and verify an existing manifest.
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

const verifyManifestVersion = "v0.2"
const verifyManifestWebsite = "https://github.com/robert-wallis/VerifyManifest"

type commandFlag struct {
	RootDir          string
	ManifestFilename string
	UnknownFilename  string
	infoLog          *log.Logger
	errorLog         *log.Logger
	exit             func(code int)
}

var gFlags = commandFlag{}

func init() {
	flag.StringVar(&gFlags.RootDir, "root", ".", "Root folder to calculate Sum.")
	flag.StringVar(&gFlags.ManifestFilename, "manifest", "manifest.json", "Manifest file name.")
	flag.StringVar(&gFlags.UnknownFilename, "unknown", "", "A text manifest file that contains hash sums in an unknown format.  Every sum in \"unknown\" file must be present in directory to pass.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\nVersion %s\n%s\n\n", os.Args[0], verifyManifestVersion, verifyManifestWebsite)
		flag.PrintDefaults()
	}
	gFlags.infoLog = log.New(os.Stdout, "", 0)
	gFlags.errorLog = log.New(os.Stderr, "", 0)
	gFlags.exit = os.Exit
}

func main() {
	flag.Parse()
	hasher := NewFolderHasher(gFlags.ManifestFilename, gFlags.UnknownFilename, gFlags.infoLog, gFlags.errorLog)
	err := hasher.HashFolder(gFlags.RootDir)
	if err != nil {
		gFlags.errorLog.Print(err)
		gFlags.exit(1)
	}
}
