package main

import (
	"os"
	"testing"
)

func Test_Filter_filterFiles(t *testing.T) {
	// GIVEN the test folder
	dirname := "test_data"
	manifestFilename := "manifest.json"
	dir, _ := os.Open(dirname)
	files, _ := dir.Readdir(0)

	// WHEN the files are filtered for the folder
	c := make(chan string)
	go filterFiles(files, manifestFilename, c)

	// THEN the first file should be a.txt
	a := <-c
	if a != "a.txt" {
		t.Fatalf("Expected a.txt got %v", a)
	}

	b := <-c
	if b != "b.txt" {
		t.Fatalf("Expected b.txt got %v", b)
	}

	end := <-c
	if end != "" {
		t.Fatalf("No more files should have been found, but found %v", end)
	}
}
