// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import (
	"os"
	"fmt"
	"encoding/json"
	"path"
)

type Manifest map[string]Sum

func (m *Manifest) Load(dirName, manifestName string) error {
	filename := path.Join(dirName, manifestName)
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Couldn't open manifest %v: %v", filename, err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	manifest := map[string]Sum{}
	err = dec.Decode(&manifest)
	if err != nil {
		return fmt.Errorf("Couldn't understand manifest file format %v: %v", filename, err)
	}
	return nil
}

func (m *Manifest) Save(dirName string, manifestName string) error {
	filename := path.Join(dirName, manifestName)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Couldn't create manifest file %v: %v", filename, err)
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	err = enc.Encode(m)
	if err != nil {
		return err
	}
	return nil
}