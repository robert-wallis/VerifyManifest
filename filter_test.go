// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

package main

import (
	"os"
	"path"
	"testing"
)

func Test_Filter_filterFiles(t *testing.T) {
	// GIVEN the test folder
	dirName := "test_data"
	manifestFilename := "manifest.json"
	dir, _ := os.Open(dirName)
	fileChan := make(chan *pathFileInfo)
	go func() {
		defer close(fileChan)
		files, _ := dir.Readdir(0)
		for f := range files {
			fileChan <- &pathFileInfo{
				files[f],
				path.Join(dirName, files[f].Name()),
				files[f].Name(),
			}
		}
	}()

	done := make(chan struct{})
	out := make(chan *pathFileInfo)

	// WHEN the files are filtered for the folder
	go filterFiles(done, fileChan, manifestFilename, out)

	// THEN the first file should be a.txt
	a := <-out
	if a.Name() != "a.txt" {
		t.Fatalf("Expected a.txt got %v", a)
	}

	b := <-out
	if b.Name() != "b.txt" {
		t.Fatalf("Expected b.txt got %v", b)
	}

	end := <-out
	if end != nil {
		t.Fatalf("No more files should have been found, but found %v", end.Name())
	}
}
