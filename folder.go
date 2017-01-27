package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type fileNameSum struct {
	FileName string
	Sum      Sum
}

func streamHashes(dirName string, files chan string, result chan *fileNameSum) error {
	defer close(result)
	for f := range files {
		fileName := path.Join(dirName, f)
		fs := &fileNameSum{
			FileName: f,
		}
		err := fs.Sum.Calculate(fileName)
		if err != nil {
			return err
		}
		result <- fs
	}
	return nil
}

func hashFolder(dirname, manifestFile string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	oldManifest := Manifest{}
	err = oldManifest.Load(dirname, manifestFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		fmt.Println("Continuing.")
	}
	newManifest := Manifest{}

	fileNameSums := make(chan *fileNameSum)
	filtered := make(chan string)
	go filterFiles(files, manifestFile, filtered)
	go func() {
		err := streamHashes(dirname, filtered, fileNameSums)
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
		err = newManifest.Save(dirname, manifestFile)
		if err != nil {
			return fmt.Errorf("Error saving manifest %v", err)
		}
		fmt.Printf("Saved manifest to %v\n", path.Join(dirname, manifestFile))
	} else {
		return fmt.Errorf("Some hashes failed, manifest not updated. %v", checkFailure)
	}
	return nil
}
