// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"errors"
)

type fileNameSum struct {
	FileName string
	Sum      Sum
}

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

func hashFolder(dirName, manifestFileName, unknownFileName string, infoLog, errorLog *log.Logger) error {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	newManifest := Manifest{}
	oldManifest := Manifest{}
	if len(manifestFileName) > 0 {
		if err := oldManifest.Load(dirName, manifestFileName); err != nil {
			infoLog.Println("Warning:", err)
			infoLog.Println("Continuing.")
		}
	}
	var unknownHashes *UnknownHashes
	if unknownFileName != "" {
		unknownHashes, err = LoadUnknownHashes(unknownFileName)
		if err != nil {
			return fmt.Errorf("Unable load \"unknown\" hash file: %v", err)
		}
	}

	fileNameSums := make(chan *fileNameSum)
	filtered := make(chan string)
	go filterFiles(files, manifestFileName, filtered)
	go func() {
		if err := streamHashes(dirName, filtered, fileNameSums); err != nil {
			errorLog.Println(err)
		}
	}()

	checkFailure := false
	for f := range fileNameSums {
		newManifest[f.FileName] = f.Sum
		if err := oldManifest.Verify(f.FileName, f.Sum); err != nil {
			checkFailure = true
			errorLog.Printf("Error %v: %v\n", f.FileName, err)
		}
		if unknownHashes != nil {
			unknownHashes.RemoveSum(f.Sum)
		}
		infoLog.Printf("%v\tmd5:%v\tsha1:%v\n", f.FileName, f.Sum.MD5, f.Sum.SHA1)
	}
	if unknownHashes != nil {
		for k, v := range *unknownHashes {
			checkFailure = true
			errorLog.Printf("Hash %v was in %v line %d, but not found in dir: %v", k, unknownFileName, v.LineNumber, v.Line)
		}
	}

	if !checkFailure {
		if len(manifestFileName) > 0 {
			if err := newManifest.Save(dirName, manifestFileName); err != nil {
				return fmt.Errorf("Error saving manifest %v", err)
			}
			infoLog.Printf("Saved manifest to %v\n", path.Join(dirName, manifestFileName))
		}
	} else {
		return errors.New("Some hashes failed, manifest not updated.")
	}
	return nil
}
