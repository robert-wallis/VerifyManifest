// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var rootDir = flag.String("root", ".", "Root folder to calculate Sum.")
var manifestFilename = flag.String("manifest", "manifest.json", "Manifest file name.")

func main() {
	flag.Parse()

	files, err := ioutil.ReadDir(*rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	oldManifest := Manifest{}
	err = oldManifest.Load(*rootDir, *manifestFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		fmt.Println("Continuing.")
	}
	newManifest := Manifest{}

	fileNameSums := make(chan *fileNameSum)
	go func() {
		err := streamHashes(*rootDir, filterFiles(files), fileNameSums)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}()

	checkFailure := false
	for f := range fileNameSums {
		if oldManifest != nil {
			old, ok := oldManifest[f.FileName]
			if ok {
				if err := old.Verify(f.Sum); err != nil {
					checkFailure = true
					fmt.Fprintf(os.Stderr, "Error %v: %v\n", f.FileName, err)
				}
			}
		}
		newManifest[f.FileName] = f.Sum
		fmt.Printf("%v\tmd5:%v\tsha1:%v\n", f.FileName, f.Sum.MD5, f.Sum.SHA1)
	}

	if !checkFailure {
		err = newManifest.Save(*rootDir, *manifestFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving manifest %v\n", err)
		}
		fmt.Printf("Saved manifest to %v\n", path.Join(*rootDir, *manifestFilename))
	} else {
		fmt.Fprintln(os.Stderr, "Some hashes failed, manifest not updated.")
	}
}

type fileNameSum struct {
	FileName string
	Sum      Sum
}

func streamHashes(dirName string, files []string, result chan *fileNameSum) error {
	for f := range files {
		fileName := path.Join(dirName, files[f])
		fs := &fileNameSum{
			FileName: files[f],
		}
		err := fs.Sum.Calculate(fileName)
		if err != nil {
			return err
		}
		result <- fs
	}
	close(result)
	return nil
}

func filterFiles(files []os.FileInfo) []string {
	result := []string{}
	for f := range files {
		if files[f].IsDir() {
			continue
		}
		if files[f].Name() == *manifestFilename {
			continue
		}
		result = append(result, files[f].Name())
	}
	return result
}
