// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"
)

func Test_streamHashes(t *testing.T) {
	// GIVEN the test folder and a stream of files
	dirName := "test_data"
	files := make(chan string)
	go func() {
		files <- "a.txt"
		close(files)
	}()

	// WHEN the files are streamed
	results := make(chan *fileNameSum)
	go streamHashes(dirName, files, results)

	// THEN the hashes are generated for a
	a := <-results
	if a.Sum.MD5 != "0cc175b9c0f1b6a831c399e269772661" {
		t.Fatalf("a.txt hash unexpected %v", a)
	}

	// THEN the channel is closed
	end := <-results
	if end != nil {
		t.Fatalf("Should have ended the stream %v", end)
	}
}

func Test_streamHashes_FileError(t *testing.T) {
	// GIVEN the test folder, but a file that doesn't
	dirName := "test_data"
	files := make(chan string)
	go func() {
		files <- "noexist"
		files <- "b.txt"
		close(files)
	}()

	// WHEN the files are streamed
	results := make(chan *fileNameSum)
	err := streamHashes(dirName, files, results)

	// THEN an error should have occurred
	if err == nil {
		t.Error("Stream should have failed, because file didn't exist.")
	}

	for f := range results {
		t.Errorf("No results, but got %v", f)
	}
}

func Test_hashFolder_BadManifest(t *testing.T) {
	// GIVEN the test folder
	dirName := "test_data"
	manifestFile := "bad_manifests/bad_b.json"
	var infoBuffer bytes.Buffer
	var errorBuffer bytes.Buffer

	// WHEN the folder is hashed
	infoLog := log.New(&infoBuffer, "", 0)
	errorLog := log.New(&errorBuffer, "", 0)
	err := hashFolder(dirName, manifestFile, "", infoLog, errorLog)
	if err == nil {
		t.Error("Should have returned a failure.")
	}

	// THEN there should be an error log that the manifest failed for b.txt
	s := errorBuffer.String()
	if len(s) == 0 {
		t.Fatal("Should have had something in the error log, but didnt.")
	}
}

func Test_hashFolder_DirError(t *testing.T) {
	// GIVEN a missing folder
	dirName := "noexist"
	manifestName := "manifest.json"
	var infoBuffer bytes.Buffer
	var errorBuffer bytes.Buffer

	// WHEN the folder is hashed
	infoLog := log.New(&infoBuffer, "", 0)
	errorLog := log.New(&errorBuffer, "", 0)
	err := hashFolder(dirName, manifestName, "", infoLog, errorLog)

	// THEN there should be an error log that the folder didn't exist
	if err == nil {
		t.Error("Should have returned a failure.")
	}
}

func Test_hashFolder_FileError(t *testing.T) {
	// GIVEN a missing file
	dirName := "test_data"
	manifestName := "noexist"
	var infoBuffer bytes.Buffer
	var errorBuffer bytes.Buffer
	defer os.Remove(path.Join(dirName, manifestName))

	// WHEN the folder is hashed
	infoLog := log.New(&infoBuffer, "", 0)
	errorLog := log.New(&errorBuffer, "", 0)
	err := hashFolder(dirName, manifestName, "", infoLog, errorLog)

	// THEN there should be not be an error that the manifest didn't exist
	if err != nil {
		t.Errorf("Should have returned a failure. %v", err)
	}

	// THEN there should have been some info in the log
	is := infoBuffer.String()
	if len(is) == 0 || !strings.Contains(is, "Warning") {
		t.Errorf("Missing manifest should have a warning: %v", is)
	}

	// THEN there shouldn't be any data in the error folder
	es := strings.Trim(errorBuffer.String(), string([]byte{0}))
	if len(es) != 0 {
		t.Errorf("There shouldn't be an error if the file didn't exist. %v %v", es, len(es))
	}
}

func Test_hashFolder_NoManifest(t *testing.T) {
	// GIVEN no manifest filename (none wanted)
	dirName := "test_data"
	infoBuffer := bytes.Buffer{}
	infoLog := log.New(&infoBuffer, "", 0)
	errorLog := log.New(&bytes.Buffer{}, "", 0)

	// THEN it shouldn't error
	err := hashFolder(dirName, "", "", infoLog, errorLog)
	if err != nil {
		t.Errorf("Missing manifest file should be ok, just don't save and print results. %v", err)
	}

	is := infoBuffer.String()
	aMd5 := "0cc175b9c0f1b6a831c399e269772661"
	if !strings.Contains(is, aMd5) {
		t.Errorf("%v was not in the info buffer, but should have been: %v", aMd5, is)
	}
}

func Test_hashFolder_FileSaveError(t *testing.T) {
	// GIVEN a missing file
	dirName := "test_data"
	manifestName := "."
	var infoBuffer bytes.Buffer
	var errorBuffer bytes.Buffer

	// WHEN the folder is hashed
	infoLog := log.New(&infoBuffer, "", 0)
	errorLog := log.New(&errorBuffer, "", 0)
	err := hashFolder(dirName, manifestName, "", infoLog, errorLog)

	// THEN there should be not be an error that the manifest didn't exist
	if err == nil {
		t.Errorf("Should have returned a that it cant save manifest to a dir. %v", err)
	}
}

func Test_hashFolder_Unknown_Success(t *testing.T) {
	// GIVEN an unknown hash set with all the values in the manifest
	dirName := "test_data"
	manifestFilename := "manifest.json"
	unknownFilename := "test_data/other_manifests/powershell.md5.txt"
	infoLog := log.New(&bytes.Buffer{}, "", 0)
	var errorBuffer bytes.Buffer
	errorLog := log.New(&errorBuffer, "", 0)

	// WHEN the folder is hashed
	err := hashFolder(dirName, manifestFilename, unknownFilename, infoLog, errorLog)

	// THEN there should not be an error
	if err != nil {
		t.Error(err)
	}

	// THEN there should be no error log
	errorString := errorBuffer.String()
	if len(errorString) != 0 {
		t.Error(errorString)
	}
}

func Test_hashFolder_Unknown_MissingHash(t *testing.T) {
	// GIVEN an unknown hash set with a missing hash
	dirName := "test_data"
	manifestFilename := "manifest.json"
	unknownFilename := "test_data/bad_manifests/powershell.extra.md5.txt"
	infoLog := log.New(&bytes.Buffer{}, "", 0)
	var errorBuffer bytes.Buffer
	errorLog := log.New(&errorBuffer, "", 0)

	// WHEN the folder is hashed
	err := hashFolder(dirName, manifestFilename, unknownFilename, infoLog, errorLog)

	// THEN there should be an error
	if err == nil {
		t.Error("Should have been an error about the missing hash.")
	}

	// THEN there should be an error for the specific missing hash
	missingHash := "48E2A9E44A8D96A6B07EAB35A86AA556"
	errorString := errorBuffer.String()
	if len(errorString) == 0 {
		t.Fatal("Error log was empty, should have an error for the hash.")
	}
	if !strings.Contains(errorString, missingHash) {
		t.Errorf("Error log didn't contain missing hash %v: %v", missingHash, errorString)
	}
}

func Test_hashFolder_Unknown_FileLoadError(t *testing.T) {
	// GIVEN a missing unknown file
	dirName := "test_data"
	unknownFilename := "noexist"
	var infoBuffer bytes.Buffer
	var errorBuffer bytes.Buffer

	// WHEN the folder is hashed
	infoLog := log.New(&infoBuffer, "", 0)
	errorLog := log.New(&errorBuffer, "", 0)
	err := hashFolder(dirName, "", unknownFilename, infoLog, errorLog)

	// THEN there should be a failure that the file didn't exist
	if err == nil {
		t.Errorf("Should have returned a failure. %v", err)
	}
	if !strings.Contains(fmt.Sprintf("%v", err), unknownFilename) {
		t.Errorf("Error should contain filename %v: %v", unknownFilename, err)
	}
}
