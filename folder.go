// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

type fileNameSum struct {
	FileName string
	Sum      Sum
}

type folderHasher struct {
	errorLog         *log.Logger
	infoLog          *log.Logger
	manifestFileName string
	unknownFileName  string
}

type pathFileInfo struct {
	os.FileInfo
	path string
	name string
}

// NewFolderHasher returns a folderHasher that can be used to verify the contents of every file in the folder.
// `manifestFileName` is the file in the folder that contains the manifest, ex. `manifest.json`
// `unknownFileName` is a text file that contains some random manifest format with hashes in it.  For example `openssl md5 * > manifest.txt`
// `infoLog` the logger where information messages are printed, i.e. os.Stdout
// `errorLog` the logger where errors are printed, i.e. os.Stderr
func NewFolderHasher(manifestFileName, unknownFileName string, infoLog, errorLog *log.Logger) *folderHasher {
	return &folderHasher{
		errorLog:         errorLog,
		infoLog:          infoLog,
		manifestFileName: manifestFileName,
		unknownFileName:  unknownFileName,
	}
}

// HashFolder goes through the directory, calculate all the hashes, and save them to a manifest.
func (h *folderHasher) HashFolder(dirName string) error {
	oldManifest, unknownHashes, err := h.loadPreviousHashes(dirName)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	files := make(chan *pathFileInfo)
	filteredFiles := make(chan *pathFileInfo)
	fileNameSums := make(chan *fileNameSum)
	go walkFolder(dirName, done, files)
	go filterFiles(done, files, h.manifestFileName, filteredFiles)
	go func() {
		if err := streamHashes(done, filteredFiles, fileNameSums); err != nil {
			h.errorLog.Println(err)
		}
	}()
	newManifest, verifyFail := h.verifyFiles(done, fileNameSums, oldManifest, unknownHashes)
	verifyUnknownFail := h.verifyUnknownHashes(unknownHashes)
	if verifyFail || verifyUnknownFail {
		return errors.New("Some hashes failed, manifest not updated.")
	}

	if len(h.manifestFileName) > 0 {
		if err := newManifest.Save(dirName, h.manifestFileName); err != nil {
			return fmt.Errorf("Error saving manifest %v", err)
		}
		h.infoLog.Printf("Saved manifest to %v\n", path.Join(dirName, h.manifestFileName))
	}
	return nil
}

// walkFolder will walk through all the files in dirName and source them into the files channel
func walkFolder(dirName string, done chan struct{}, files chan *pathFileInfo) (err error) {
	defer close(files)
	err = filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			close(done)
			return err
		}
		if path == dirName {
			return nil
		}
		files <- &pathFileInfo{
			FileInfo: info,
			path:     path,
			name:     path[len(dirName)+1:],
		}
		return nil
	})
	return
}

// any unknown hashes left are failures
func (h *folderHasher) verifyUnknownHashes(unknownHashes *UnknownHashes) bool {
	verifyFail := false
	if unknownHashes != nil {
		for k, v := range *unknownHashes {
			verifyFail = true
			h.errorLog.Printf("Hash %v was in %v line %d, but not found in dir: %v", k, h.unknownFileName, v.LineNumber, v.Line)
		}
	}
	return verifyFail
}

// go through all the files in files stream, calculate the hash, and then send the result over the result stream
func streamHashes(done chan struct{}, files chan *pathFileInfo, result chan *fileNameSum) error {
	defer close(result)
	for file := range files {
		fs := &fileNameSum{
			FileName: file.name,
		}
		err := fs.Sum.Calculate(file.path)
		if err != nil {
			close(done)
			return err
		}
		select {
		case <-done:
			return nil
		case result <- fs:
		}
	}
	return nil
}

// go though all the hashes in the fileNameSums stream, save them in the newManifest, and remove them from unknownHashes
func (h *folderHasher) verifyFiles(done chan struct{}, fileNameSums chan *fileNameSum, oldManifest *Manifest, unknownHashes *UnknownHashes) (newManifest *Manifest, verifyFail bool) {
	newManifest = &Manifest{}
	for f := range fileNameSums {
		(*newManifest)[f.FileName] = f.Sum
		if err := oldManifest.Verify(f.FileName, f.Sum); err != nil {
			verifyFail = true
			h.errorLog.Printf("Error %v: %v\n", f.FileName, err)
		}
		if unknownHashes != nil {
			unknownHashes.RemoveSum(f.Sum)
		}
		h.infoLog.Printf("%v\tmd5:%v\tsha1:%v\n", f.FileName, f.Sum.MD5, f.Sum.SHA1)
		select {
		case <-done:
			return
		default:
		}
	}
	return
}

// load the oldManifest and/or an unknownHahses file
func (h *folderHasher) loadPreviousHashes(dirName string) (oldManifest *Manifest, unknownHashes *UnknownHashes, err error) {
	oldManifest = &Manifest{}

	if len(h.manifestFileName) > 0 {
		if err := oldManifest.Load(dirName, h.manifestFileName); err != nil {
			h.infoLog.Println("Warning:", err)
			h.infoLog.Println("Continuing.")
		}
	}
	if h.unknownFileName != "" {
		unknownHashes, err = LoadUnknownHashes(h.unknownFileName)
		if err != nil {
			return nil, nil, fmt.Errorf("Unable load \"unknown\" hash file: %v", err)
		}
	}
	return
}
