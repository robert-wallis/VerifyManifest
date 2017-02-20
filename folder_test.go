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
	files := make(chan *pathFileInfo)
	go func() {
		fi, _ := os.Stat("test_data/a.txt")
		files <- &pathFileInfo{fi, path.Join(dirName, fi.Name()), fi.Name()}
		close(files)
	}()
	done := make(chan struct{})

	// WHEN the files are streamed
	results := make(chan *fileNameSum)
	go func() {
		err := streamHashes(done, files, results)

		// THEN there were no errors
		if err != nil {
			t.Fatal(err)
		}
	}()

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
	// GIVEN an error caused the done to be called
	dirName := "test_data"
	files := make(chan *pathFileInfo)
	done := make(chan struct{})
	go func() {
		close(done)
		fi, _ := os.Stat("test_data/a.txt")
		files <- &pathFileInfo{fi, path.Join(dirName, fi.Name()), fi.Name()}
		fi, _ = os.Stat("test_data/b.txt")
		files <- &pathFileInfo{fi, path.Join(dirName, fi.Name()), fi.Name()}
		close(files)
	}()

	// WHEN the files are streamed
	results := make(chan *fileNameSum)
	err := streamHashes(done, files, results)

	// THEN there shouldn't be any errors
	if err != nil {
		t.Errorf("Should be no return error if just done: %v", err)
	}

	// THEN there shouldn't be any results
	for f := range results {
		t.Errorf("No results, but got %v", f)
	}
}

func Test_loadPreviousHashes_MissingManifest(t *testing.T) {
	// GIVEN a missing manifest file
	dirName := "test_data"
	manifestName := "noexist"
	infoBuffer, errorBuffer, h := makeTestFolderHasher(manifestName, "")

	// WHEN the manifest is loaded
	_, _, err := h.loadPreviousHashes(dirName)

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

func Test_loadPreviousHashes_FileLoadError(t *testing.T) {
	// GIVEN a missing unknown file
	dirName := "test_data"
	manifestFilename := ""
	unknownFilename := "noexist"
	_, _, h := makeTestFolderHasher(manifestFilename, unknownFilename)

	// WHEN the folder is hashed
	_, _, err := h.loadPreviousHashes(dirName)

	// THEN there should be a failure that the file didn't exist
	if err == nil {
		t.Errorf("Should have returned a failure. %v", err)
	}
	if !strings.Contains(fmt.Sprintf("%v", err), unknownFilename) {
		t.Errorf("Error should contain filename %v: %v", unknownFilename, err)
	}
}

func Test_hashFolder_DirError(t *testing.T) {
	// GIVEN a missing folder
	dirName := "noexist"
	manifestName := "manifest.json"
	_, _, h := makeTestFolderHasher(manifestName, "")

	// WHEN the folder is hashed
	err := h.HashFolder(dirName)

	// THEN there should be an error log that the folder didn't exist
	if err == nil {
		t.Error("Should have returned a failure.")
	}
}

func Test_hashFolder_NoManifest(t *testing.T) {
	// GIVEN no manifest filename (none wanted)
	dirName := "test_data"
	manifestFile := ""
	infoBuffer, errorBuffer, h := makeTestFolderHasher(manifestFile, "")

	// THEN it shouldn't error
	err := h.HashFolder(dirName)
	if err != nil {
		t.Errorf("Missing manifest file should be ok, just don't save and print results. %v", err)
	}
	if errorBuffer.Len() > 0 {
		t.Errorf("The error buffer should be empty: %v", errorBuffer)
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
	_, _, h := makeTestFolderHasher(manifestName, "")

	// WHEN the folder is hashed
	err := h.HashFolder(dirName)

	// THEN there should be not be an error that the manifest didn't exist
	if err == nil {
		t.Errorf("Should have returned a that it cant save manifest to a dir. %v", err)
	}
}

func Test_hashFolder_ManifestWithInvalidHash(t *testing.T) {
	// GIVEN a manifest with incorrect hashes
	dirName := "test_data"
	manifestFile := "bad_manifests/bad_b.json"
	_, errorBuffer, h := makeTestFolderHasher(manifestFile, "")

	// WHEN the manifest is loaded
	err := h.HashFolder(dirName)
	if err == nil {
		t.Error("Should have returned a failure.")
	}

	// THEN there should be an error log that the manifest failed for b.txt
	s := errorBuffer.String()
	if len(s) == 0 {
		t.Fatal("Should have had something in the error log, but didnt.")
	}
}

func Test_hashFolder_Unknown_Success(t *testing.T) {
	// GIVEN an unknown hash set with all the values in the manifest
	dirName := "test_data"
	manifestFilename := "manifest.json"
	unknownFilename := "test_data/other_manifests/powershell.md5.txt"
	_, errorBuffer, h := makeTestFolderHasher(manifestFilename, unknownFilename)

	// WHEN the folder is hashed
	err := h.HashFolder(dirName)

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
	_, errorBuffer, h := makeTestFolderHasher(manifestFilename, unknownFilename)

	// WHEN the folder is hashed
	err := h.HashFolder(dirName)

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

func makeTestFolderHasher(manifestFileName, unknownFileName string) (infoBuffer, errorBuffer *bytes.Buffer, hasher *folderHasher) {
	infoBuffer = &bytes.Buffer{}
	errorBuffer = &bytes.Buffer{}
	infoLog := log.New(infoBuffer, "", 0)
	errorLog := log.New(errorBuffer, "", 0)
	hasher = NewFolderHasher(manifestFileName, unknownFileName, infoLog, errorLog)
	return
}
