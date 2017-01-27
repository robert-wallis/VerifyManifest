package main

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

func Test_LoadUnknown_BadFile(t *testing.T) {
	// GIVEN a bad filename
	fileName := "noexist"

	// WHEN trying to open it
	_, err := LoadUnknownForHashes(fileName)

	// THEN it should error
	if err == nil {
		t.Errorf("Should have errored about the file not existing. %v", err)
	}
}

func Test_LoadUnknown_Powershell_MD5(t *testing.T) {
	// GIVEN a powershell hashsum file
	// generated with `get-filehash -algorithm md5 *.txt > powershell.md5.txt`
	fileName := "test_data/bad_manifests/powershell.md5.txt"

	// WHEN it is scanned for MD5 files
	results, err := LoadUnknownForHashes(fileName)

	// THEN it should contain the expected hashes in the right locations
	if err != nil {
		t.Fatal(err)
	}
	aSum := "0cc175b9c0f1b6a831c399e269772661"
	v, ok := results[aSum]
	if !ok {
		t.Error("Didn't find the hash for a.txt")
	} else if v.LineNumber != 3 {
		t.Errorf("a.txt found on the wrong line number %v", v.LineNumber)
	}

	bSum := "92eb5ffee6ae2fec3ad71c777531578f"
	v, ok = results[bSum]
	if !ok {
		t.Error("Didn't find the hash for b.txt")
	} else if v.LineNumber != 4 {
		t.Errorf("b.txt found on the wrong line number %v", v.LineNumber)
	}
}

func Test_LoadUnknown_Powershell_SHA1(t *testing.T) {
	// GIVEN a powershell hashsum file
	// generated with `get-filehash -algorithm SHA1 *.txt > powershell.sha1.txt`
	fileName := "test_data/bad_manifests/powershell.sha1.txt"

	// WHEN it is scanned for SHA1 files
	results, err := LoadUnknownForHashes(fileName)

	// THEN it should contain the expected hashes in the right locations
	if err != nil {
		t.Fatal(err)
	}
	aSum := "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8"
	v, ok := results[aSum]
	if !ok {
		t.Error("Didn't find the hash for a.txt")
	} else if v.LineNumber != 3 {
		t.Errorf("a.txt found on the wrong line number %v", v.LineNumber)
	}

	bSum := "e9d71f5ee7c92d6dc9e92ffdad17b8bd49418f98"
	v, ok = results[bSum]
	if !ok {
		t.Error("Didn't find the hash for b.txt")
	} else if v.LineNumber != 4 {
		t.Errorf("b.txt found on the wrong line number %v", v.LineNumber)
	}
}
