package main

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func Test_loadUnknown_BadFile(t *testing.T) {
	// GIVEN a bad filename
	fileName := "noexist"

	// WHEN trying to open it
	_, err := LoadUnknownForHashes(fileName)

	// THEN it should error
	if err == nil {
		t.Errorf("Should have errored about the file not existing. %v", err)
	}
}

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
	if val := MD5InString(valid); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the md5 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the md5 sum.")
		}
	}

	// WHEN it's at the beginning of the string
	// THEN it should find it
	beginning := valid + " something else"
	if val := MD5InString(beginning); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the md5 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the md5 sum.")
		}
	}

	// WHEN it's at the end of the string
	// THEN it should find it
	end := "something else " + valid
	if val := MD5InString(end); val == nil || *val != valid {
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
	if val := MD5InString("0123456789abcdefABCDEF012345678"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := MD5InString("a 0123456789abcdefABCDEF012345678"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := MD5InString("0123456789abcdefABCDEF012345678 a"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}

	// GIVEN a string with something too long to be an md5
	// WHEN md5String is called
	// THEN it should return that md5
	if val := MD5InString("0123456789abcdefABCDEF0123456789a"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := MD5InString("b 0123456789abcdefABCDEF0123456789a"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := MD5InString("0123456789abcdefABCDEF0123456789a b"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
}

func Test_SHA1InString(t *testing.T) {
	// GIVEN a string with something that looks like an sha1
	valid := "0123456789abcdefABCDEF0123456789abcdef40"
	// WHEN sha1String is called
	// THEN it should return that sha1
	if val := SHA1InString(valid); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the sha1 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the sha1 sum.")
		}
	}

	// WHEN it's at the beginning of the string
	// THEN it should find it
	beginning := valid + " something else"
	if val := SHA1InString(beginning); val == nil || *val != valid {
		if val != nil {
			t.Errorf("Expecting to find the sha1 sum, found \"%v\"", *val)
		} else {
			t.Error("Didn't find the sha1 sum.")
		}
	}

	// WHEN it's at the end of the string
	// THEN it should find it
	end := "something else " + valid
	if val := SHA1InString(end); val == nil || *val != valid {
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
	if val := SHA1InString("0123456789abcdefABCDEF0123456789abcdef4"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := SHA1InString("a 0123456789abcdefABCDEF0123456789abcdef4"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := SHA1InString("0123456789abcdefABCDEF0123456789abcdef4 a"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}

	// GIVEN a string with something too long to be an sha1
	// WHEN sha1String is called
	// THEN it should return that sha1
	if val := SHA1InString("0123456789abcdefABCDEF0123456789abcdef40a"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := SHA1InString("b 0123456789abcdefABCDEF0123456789abcdef40a"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
	if val := SHA1InString("0123456789abcdefABCDEF0123456789abcdef40a b"); val != nil {
		t.Errorf("Should have found nothing, found %v", *val)
	}
}