// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import "testing"

func Test_Sum_Calculate(t *testing.T) {
	// GIVEN a tested file
	filename := "test_data/a.txt"
	md5 := "0cc175b9c0f1b6a831c399e269772661"
	sha1 := "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8"

	// WHEN the sums are calculated
	sum := Sum{}
	if err := sum.Calculate(filename); err != nil {
		t.Error(err)
	}

	// THEN the values should be expected
	if md5 != sum.MD5 {
		t.Errorf("Exected MD5 %v got %v", md5, sum.MD5)
	}

	if sha1 != sum.SHA1 {
		t.Errorf("Exected SHA1 %v got %v", md5, sum.SHA1)
	}
}

func Test_Sum_Calculate_FileError(t *testing.T) {
	// GIVEN a file that doesn't exist
	filename := "test_data/noexist"

	// WHEN the ums are calculated
	sum := Sum{}
	err := sum.Calculate(filename)

	// THEN an error should have been generated
	if err == nil {
		t.Errorf("Expecting an error but didn't get one for %v", filename)
	}
}

func Test_Sum_Verify(t *testing.T) {
	// GIVEN a tested file
	sum := Sum{
		MD5:  "md5",
		SHA1: "sha1",
	}

	// WHEN the sum is tested against the same values
	correct := Sum{
		MD5:  "md5",
		SHA1: "sha1",
	}

	// THEN the sum should be correct
	if err := sum.Verify(correct); err != nil {
		t.Errorf("Expected verify to work, but didn't. %v", err)
	}

	// WHEN the sum is tested against a wrong md5
	bad_md5 := Sum{
		MD5:  "x",
		SHA1: "sha1",
	}

	// THEN it should fail
	if err := sum.Verify(bad_md5); err == nil {
		t.Errorf("Expecting failure with MD5 but didn't error. %v", bad_md5)
	}

	// WHEN the sum is tested against the wrong sha1
	bad_sha1 := Sum{
		MD5:  "md5",
		SHA1: "x",
	}

	// THEN it should fail
	if err := sum.Verify(bad_sha1); err == nil {
		t.Errorf("Expecting failure with SHA1 but didn't error. %v", bad_sha1)
	}
}
