package main

import (
	"bytes"
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
	infoBuf := make([]byte, 1000000)
	infoBuffer := bytes.NewBuffer(infoBuf)
	errorBuf := make([]byte, 1000000)
	errorBuffer := bytes.NewBuffer(errorBuf)

	// WHEN the folder is hashed
	infoLog := log.New(infoBuffer, "", 0)
	errorLog := log.New(errorBuffer, "", 0)
	err := hashFolder(dirName, manifestFile, infoLog, errorLog)
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
	infoBuf := make([]byte, 1000000)
	infoBuffer := bytes.NewBuffer(infoBuf)
	errorBuf := make([]byte, 1000000)
	errorBuffer := bytes.NewBuffer(errorBuf)

	// WHEN the folder is hashed
	infoLog := log.New(infoBuffer, "", 0)
	errorLog := log.New(errorBuffer, "", 0)
	err := hashFolder(dirName, manifestName, infoLog, errorLog)

	// THEN there should be an error log that the folder didn't exist
	if err == nil {
		t.Error("Should have returned a failure.")
	}
}

func Test_hashFolder_FileError(t *testing.T) {
	// GIVEN a missing file
	dirName := "test_data"
	manifestName := "noexist"
	infoBuf := make([]byte, 1000000)
	infoBuffer := bytes.NewBuffer(infoBuf)
	errorBuf := make([]byte, 1000000)
	errorBuffer := bytes.NewBuffer(errorBuf)
	defer os.Remove(path.Join(dirName, manifestName))

	// WHEN the folder is hashed
	infoLog := log.New(infoBuffer, "", 0)
	errorLog := log.New(errorBuffer, "", 0)
	err := hashFolder(dirName, manifestName, infoLog, errorLog)

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

func Test_hashFolder_FileSaveError(t *testing.T) {
	// GIVEN a missing file
	dirName := "test_data"
	manifestName := "."
	infoBuf := make([]byte, 1000000)
	infoBuffer := bytes.NewBuffer(infoBuf)
	errorBuf := make([]byte, 1000000)
	errorBuffer := bytes.NewBuffer(errorBuf)

	// WHEN the folder is hashed
	infoLog := log.New(infoBuffer, "", 0)
	errorLog := log.New(errorBuffer, "", 0)
	err := hashFolder(dirName, manifestName, infoLog, errorLog)

	// THEN there should be not be an error that the manifest didn't exist
	if err == nil {
		t.Errorf("Should have returned a that it cant save manifest to a dir. %v", err)
	}
}
