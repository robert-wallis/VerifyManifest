// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
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

func hashFolder(dirName, manifestFile string, infoLog, errorLog *log.Logger) error {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	oldManifest := Manifest{}
	if err := oldManifest.Load(dirName, manifestFile); err != nil {
		infoLog.Println("Warning:", err)
		infoLog.Println("Continuing.")
	}
	newManifest := Manifest{}

	fileNameSums := make(chan *fileNameSum)
	filtered := make(chan string)
	go filterFiles(files, manifestFile, filtered)
	go func() {
		if err := streamHashes(dirName, filtered, fileNameSums); err != nil {
			errorLog.Println(err)
		}
	}()

	checkFailure := false
	for f := range fileNameSums {
		if oldManifest != nil {
			if err := oldManifest.Verify(f.FileName, f.Sum); err != nil {
				checkFailure = true
				errorLog.Printf("Error %v: %v\n", f.FileName, err)
			}
		}
		newManifest[f.FileName] = f.Sum
		infoLog.Printf("%v\tmd5:%v\tsha1:%v\n", f.FileName, f.Sum.MD5, f.Sum.SHA1)
	}

	if !checkFailure {
		if err := newManifest.Save(dirName, manifestFile); err != nil {
			return fmt.Errorf("Error saving manifest %v", err)
		}
		infoLog.Printf("Saved manifest to %v\n", path.Join(dirName, manifestFile))
	} else {
		return fmt.Errorf("Some hashes failed, manifest not updated. %v", checkFailure)
	}
	return nil
}
