// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

package manifest

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func Test_hexRune(t *testing.T) {
	// GIVEN a list of valid chars
	valid := "0123456789abcdefABCDEF"

	// WHEN checked if it's a rune
	reader := bufio.NewReader(strings.NewReader(valid))
	for {
		r, size, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			break
		}
		if size <= 0 {
			t.Errorf("Size <= 0, %v", size)
			break
		}
		// THEN it should be true
		if ok := hexRune(r); !ok {
			t.Errorf("Rune '%v' was considered invalid.", string(r))
		}
	}
}

func Test_MD5InString(t *testing.T) {
	// GIVEN a string with something that looks like an MD5
	valid := "0123456789abcdefABCDEF0123456789"
	// WHEN md5String is called
	// THEN it should return that md5
	if val := MD5InString([]rune(valid)); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the md5 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the md5 sum.")
		}
	}

	// WHEN it's at the beginning of the string
	// THEN it should find it
	beginning := valid + " something else"
	if val := MD5InString([]rune(beginning)); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the md5 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the md5 sum.")
		}
	}

	// WHEN it's at the end of the string
	// THEN it should find it
	end := "something else " + valid
	if val := MD5InString([]rune(end)); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the md5 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the md5 sum.")
		}
	}

}

func Test_MD5InString_Unicode(t *testing.T) {
	u := "0CC175B9C0F1B6A831C399E269772661"
	if val := MD5InString([]rune(u)); val == nil || *val != u {
		if val != nil {
			t.Errorf("Expecting to find the md5 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the md5 sum.")
		}
	}
}

func Test_MD5InString_missing(t *testing.T) {
	// GIVEN a string with something that is too short to be an md5
	// WHEN md5String is called
	// THEN it should return that md5
	if val := MD5InString([]rune("0123456789abcdefABCDEF012345678")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := MD5InString([]rune("a 0123456789abcdefABCDEF012345678")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := MD5InString([]rune("0123456789abcdefABCDEF012345678 a")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}

	// GIVEN a string with something too long to be an md5
	// WHEN md5String is called
	// THEN it should return that md5
	if val := MD5InString([]rune("0123456789abcdefABCDEF0123456789a")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := MD5InString([]rune("b 0123456789abcdefABCDEF0123456789a")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := MD5InString([]rune("0123456789abcdefABCDEF0123456789a b")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
}

func Test_SHA1InString(t *testing.T) {
	// GIVEN a string with something that looks like an sha1
	valid := "0123456789abcdefABCDEF0123456789abcdef40"
	// WHEN sha1String is called
	// THEN it should return that sha1
	if val := SHA1InString([]rune(valid)); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the sha1 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the sha1 sum.")
		}
	}

	// WHEN it's at the beginning of the string
	// THEN it should find it
	beginning := valid + " something else"
	if val := SHA1InString([]rune(beginning)); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the sha1 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the sha1 sum.")
		}
	}

	// WHEN it's at the end of the string
	// THEN it should find it
	end := "something else " + valid
	if val := SHA1InString([]rune(end)); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the sha1 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the sha1 sum.")
		}
	}
}

func Test_SHA1InString_missing(t *testing.T) {
	// GIVEN a string with something that is too short to be an sha1
	// WHEN sha1String is called
	// THEN it should return that sha1
	if val := SHA1InString([]rune("0123456789abcdefABCDEF0123456789abcdef4")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := SHA1InString([]rune("a 0123456789abcdefABCDEF0123456789abcdef4")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := SHA1InString([]rune("0123456789abcdefABCDEF0123456789abcdef4 a")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}

	// GIVEN a string with something too long to be an sha1
	// WHEN sha1String is called
	// THEN it should return that sha1
	if val := SHA1InString([]rune("0123456789abcdefABCDEF0123456789abcdef40a")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := SHA1InString([]rune("b 0123456789abcdefABCDEF0123456789abcdef40a")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := SHA1InString([]rune("0123456789abcdefABCDEF0123456789abcdef40a b")); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
}

func Test_LoadUnknownHashes_BadFile(t *testing.T) {
	// GIVEN a bad filename
	fileName := "noexist"

	// WHEN trying to open it
	_, err := LoadUnknownHashes(fileName)

	// THEN it should error
	if err == nil {
		t.Errorf("Should have errored about the file not existing. %v", err)
	}
}

func Test_LoadUnknownHashes_Powershell_MD5(t *testing.T) {
	// GIVEN a powershell hashsum file
	// generated with `get-filehash -algorithm md5 *.txt > powershell.md5.txt`
	fileName := "../test_data/other_manifests/powershell.md5.txt"

	// WHEN it is scanned for MD5 files
	results, err := LoadUnknownHashes(fileName)

	// THEN it should contain the expected hashes in the right locations
	if err != nil {
		t.Fatal(err)
	}
	aSum := "0cc175b9c0f1b6a831c399e269772661"
	v, ok := results.Get(aSum)
	if !ok {
		t.Error("Didn't find the hash for a.txt")
	} else if v.LineNumber != 4 {
		t.Errorf("a.txt found on the wrong line number %v", v.LineNumber)
	}

	bSum := "92eb5ffee6ae2fec3ad71c777531578f"
	v, ok = results.Get(bSum)
	if !ok {
		t.Error("Didn't find the hash for b.txt")
	} else if v.LineNumber != 5 {
		t.Errorf("b.txt found on the wrong line number %v", v.LineNumber)
	}
}

func Test_LoadUnknownHashes_Powershell_SHA1(t *testing.T) {
	// GIVEN a powershell hashsum file
	// generated with `get-filehash -algorithm SHA1 *.txt > powershell.sha1.txt`
	fileName := "../test_data/other_manifests/powershell.sha1.txt"

	// WHEN it is scanned for SHA1 files
	results, err := LoadUnknownHashes(fileName)

	// THEN it should contain the expected hashes in the right locations
	if err != nil {
		t.Fatal(err)
	}
	aSum := "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8"
	v, ok := results.Get(aSum)
	if !ok {
		t.Error("Didn't find the hash for a.txt")
	} else if v.LineNumber != 4 {
		t.Errorf("a.txt found on the wrong line number %v", v.LineNumber)
	}

	bSum := "e9d71f5ee7c92d6dc9e92ffdad17b8bd49418f98"
	v, ok = results.Get(bSum)
	if !ok {
		t.Error("Didn't find the hash for b.txt")
	} else if v.LineNumber != 5 {
		t.Errorf("b.txt found on the wrong line number %v", v.LineNumber)
	}
}

func Test_UnknownHashes_Get(t *testing.T) {
	// GIVEN a blank hashset
	u := UnknownHashes{}

	// WHEN a hash is not present
	// THEN it should fail
	if _, ok := u.Get("nope"); ok {
		t.Error("Hash should not be present.")
	}

	// WHEN a hash is present
	hash := "123456789"
	loc := HashLocation{
		LineNumber: 21,
		Line:       "test line",
	}
	u[hash] = loc

	// THEN it should be successful (ok)
	v, ok := u.Get(hash)
	if !ok {
		t.Error("Hash should be present")
	}
	if v.LineNumber != loc.LineNumber || v.Line != loc.Line {
		t.Error("Hash value was wrong", v)
	}
}

func Test_UnknownHashes_Set(t *testing.T) {
	// GIVEN a blank hashset
	u := UnknownHashes{}

	// WHEN set is called
	hash := "123456789"
	loc := HashLocation{
		LineNumber: 21,
		Line:       "test line",
	}
	u.Set(hash, loc)

	// THEN the value should be set
	v, ok := u[hash]
	if !ok {
		t.Error("Key not set")
	}
	if v.LineNumber != loc.LineNumber || v.Line != loc.Line {
		t.Error("Hash value was wrong", v)
	}
}

func Test_UnknownHashes_Remove(t *testing.T) {
	// GIVEN a populated hash file
	u, err := LoadUnknownHashes("../test_data/other_manifests/powershell.md5.txt")
	if err != nil {
		t.Fatal(err)
	}
	testHash := "0cc175b9c0f1b6a831c399e269772661"
	if _, ok := u.Get(testHash); !ok {
		t.Fatalf("Hash was not found in test data: %v", testHash)
	}

	// WHEN the hash is removed
	u.Remove(testHash)

	// THEN it should no longer be in the set
	if _, ok := u.Get(testHash); ok {
		t.Errorf("%v should have been removed", testHash)
	}
}

func Test_UnknownHashes_RemoveSet_MD5(t *testing.T) {
	// GIVEN a populated hash file
	u, err := LoadUnknownHashes("../test_data/other_manifests/powershell.md5.txt")
	if err != nil {
		t.Fatal(err)
	}
	testHash := "0cc175b9c0f1b6a831c399e269772661"
	if _, ok := u.Get(testHash); !ok {
		t.Fatalf("Hash was not found in test data: %v", testHash)
	}

	// WHEN an MD5 hash is in a sum to be removed
	sum := Sum{
		MD5: testHash,
	}
	u.RemoveSum(sum)

	// THEN it should no longer be in the set
	if _, ok := u.Get(testHash); ok {
		t.Errorf("%v should have been removed", testHash)
	}
}

func Test_UnknownHashes_RemoveSet_SHA1(t *testing.T) {
	// GIVEN a populated hash file
	u, err := LoadUnknownHashes("../test_data/other_manifests/powershell.sha1.txt")
	if err != nil {
		t.Fatal(err)
	}
	testHash := "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8"
	if _, ok := u.Get(testHash); !ok {
		t.Fatalf("Hash was not found in test data: %v", testHash)
	}

	// WHEN an MD5 hash is in a sum to be removed
	sum := Sum{
		SHA1: testHash,
	}
	u.RemoveSum(sum)

	// THEN it should no longer be in the set
	if _, ok := u.Get(testHash); ok {
		t.Errorf("%v should have been removed", testHash)
	}
}
