// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
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

func NewFolderHasher(manifestFileName, unknownFileName string, infoLog, errorLog *log.Logger) *folderHasher {
	return &folderHasher{
		errorLog:         errorLog,
		infoLog:          infoLog,
		manifestFileName: manifestFileName,
		unknownFileName:  unknownFileName,
	}
}

// go through the directory, calculate all the hashes, and save them to a manifest
func (h *folderHasher) HashFolder(dirName string) error {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	oldManifest, unknownHashes, err := h.loadPreviousHashes(dirName)
	if err != nil {
		return err
	}

	fileNameSums := make(chan *fileNameSum)
	filtered := make(chan string)
	go filterFiles(files, h.manifestFileName, filtered)
	go func() {
		if err := streamHashes(dirName, filtered, fileNameSums); err != nil {
			h.errorLog.Println(err)
		}
	}()

	newManifest, verifyFail := h.verifyFiles(fileNameSums, oldManifest, unknownHashes)
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
func streamHashes(dirName string, files chan string, result chan *fileNameSum) error {
	defer close(result)
	for fileName := range files {
		fullFileName := path.Join(dirName, fileName)
		fs := &fileNameSum{
			FileName: fileName,
		}
		err := fs.Sum.Calculate(fullFileName)
		if err != nil {
			return err
		}
		result <- fs
	}
	return nil
}

// go though all the hashes in the fileNameSums stream, save them in the newManifest, and remove them from unknownHashes
func (h *folderHasher) verifyFiles(fileNameSums chan *fileNameSum, oldManifest *Manifest, unknownHashes *UnknownHashes) (newManifest *Manifest, verifyFail bool) {
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
	}
	return newManifest, verifyFail
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
