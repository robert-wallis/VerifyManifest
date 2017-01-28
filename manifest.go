// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	if err = dec.Decode(m); err != nil {
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
	if err = enc.Encode(m); err != nil {
		return err
	}
	return nil
}

func (m *Manifest) Verify(fileName string, sum Sum) error {
	old, ok := (*m)[fileName]
	if ok {
		if err := old.Verify(sum); err != nil {
			return fmt.Errorf("Error %v: %v\n", fileName, err)
		}
	}
	return nil
}
