// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func Test_Manifest_Load(t *testing.T) {
	// GIVEN the test manifest file
	manifest := Manifest{}

	// WHEN it is loaded
	err := manifest.Load("test_data", "manifest.json")
	if err != nil {
		t.Fatal(err)
	}

	// THEN it should contain the "a.txt" and "b.txt" files
	expected := []string{
		"a.txt", "b.txt",
		fmt.Sprintf("%s%c%s", "bad_manifests", os.PathSeparator, "bad_b.json"),
		fmt.Sprintf("%s%c%s", "bad_manifests", os.PathSeparator, "powershell.extra.md5.txt"),
		fmt.Sprintf("%s%c%s", "other_manifests", os.PathSeparator, "powershell.md5.txt"),
		fmt.Sprintf("%s%c%s", "other_manifests", os.PathSeparator, "powershell.sha1.txt"),
	}
	count := 0
	for k := range manifest {
		var inExpected bool
		for e := range expected {
			if k == expected[e] {
				inExpected = true
				break
			}
		}
		if !inExpected {
			t.Error("Unexpected file in test folder", k)
		}
		count++
	}
	if count == 0 {
		t.Errorf("Should have been 2 files but were %v", count)
	}
}

func Test_Manifest_Load_FileError(t *testing.T) {
	// GIVEN a file that doesn't exist
	dirname := "test_data"
	filename := "noexist"

	// WHEN it is loaded
	manifest := Manifest{}
	err := manifest.Load(dirname, filename)

	// THEN an error should happen
	if err == nil {
		t.Fatalf("Expected an error with %v/%v but no error happened.", dirname, filename)
	}
}

func Test_Manifest_Save(t *testing.T) {
	// GIVEN a new manifest that was generated
	manifest := Manifest{}
	manifest["test.txt"] = Sum{
		MD5:  "md5",
		SHA1: "sha1",
	}

	// WHEN the manifest is saved
	err := manifest.Save(".", "test_manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test_manifest.json")

	// THEN it should save correctly
	expected := Manifest{}
	expectedText := "{\n\t\"test.txt\": {\n\t\t\"MD5\": \"md5\",\n\t\t\"SHA1\": \"sha1\"\n\t}\n}"
	expectedDec := json.NewDecoder(strings.NewReader(expectedText))
	if err = expectedDec.Decode(&expected); err != nil {
		t.Fatalf("Expected decode failed %v", err)
	}

	actual := Manifest{}
	actualFile, err := os.Open("test_manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	defer actualFile.Close()
	actualDec := json.NewDecoder(actualFile)
	if err = actualDec.Decode(&actual); err != nil {
		t.Fatalf("Actual decode failed %v", err)
	}

	for k, v := range expected {
		if v.MD5 != actual[k].MD5 {
			t.Errorf("%v MD5 value expected %v actual %v", k, v, actual[k].MD5)
		}
		if v.SHA1 != actual[k].SHA1 {
			t.Errorf("%v SHA1 value expected %v actual %v", k, v, actual[k].SHA1)
		}
	}
}

func Test_Manifest_Save_FileError(t *testing.T) {
	// GIVEN a new manifest that was generated
	// AND a bad directory
	dirname := "noexist"
	filename := "noexist"
	manifest := Manifest{}
	manifest["test.txt"] = Sum{
		MD5:  "md5",
		SHA1: "sha1",
	}

	// WHEN it is saved
	err := manifest.Save(dirname, filename)

	// THEN it should have an error
	if err == nil {
		t.Errorf("Should have had an error saving to a fake folder %v/%v but didn't", dirname, filename)
	}
}

func Test_Manifest_Verify(t *testing.T) {
	// GIVEN a manifest with a file
	m := Manifest{}
	fileName := "test.txt"
	sum := Sum{
		MD5:  "md5test",
		SHA1: "sha1test",
	}
	m[fileName] = sum

	// WHEN the file is not in the manifest
	err := m.Verify("noexist", sum)

	// THEN there should be no error
	if err != nil {
		t.Errorf("Expecting no error: %v", err)
	}

	// WHEN the file is in the manifest and the sum is the same
	err = m.Verify(fileName, sum)

	// THEN there should be no error
	if err != nil {
		t.Errorf("Expecting no error: %v", err)
	}

	// WHEN the file exists, but the sum is different
	badSum := Sum{
		MD5:  "badsum",
		SHA1: "badsum",
	}
	err = m.Verify(fileName, badSum)

	// THEN there should be an error
	if err == nil {
		t.Errorf("Expecing an error, got %v", err)
	}
}
