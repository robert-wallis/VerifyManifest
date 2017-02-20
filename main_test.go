package main

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func Test_Main(t *testing.T) {
	// GIVEN the test_data folder
	gFlags.RootDir = "test_data"
	infoBuffer := &bytes.Buffer{}
	errorBuffer := &bytes.Buffer{}
	gFlags.infoLog = log.New(infoBuffer, "", 0)
	gFlags.errorLog = log.New(errorBuffer, "", 0)
	gFlags.exit = func(code int) {
		// THEN it should exit 0 (success)
		if code != 0 {
			t.Errorf("Non-successful exit code %v", code)
		}
	}

	// WHEN VerifyManifest is run
	main()

	// THEN it should not error
	if errorBuffer.Len() > 0 {
		t.Errorf("There should be no errors: %v", errorBuffer.String())
	}

	// THEN there should be specific hashes in the output
	infoStr := infoBuffer.String()
	if !strings.Contains(infoStr, "0cc175b9c0f1b6a831c399e269772661") {
		t.Error("The hash for a.txt was not in the output.")
	}
	if !strings.Contains(infoStr, "e9d71f5ee7c92d6dc9e92ffdad17b8bd49418f98") {
		t.Error("The hash for b.txt was not in the output.")
	}
}

func Test_Main_Error(t *testing.T) {
	// GIVEN HashFolder has an error
	gFlags.RootDir = "test_data"
	gFlags.UnknownFilename = "noexist" // this should generate an error
	infoBuffer := &bytes.Buffer{}
	errorBuffer := &bytes.Buffer{}
	gFlags.infoLog = log.New(infoBuffer, "", 0)
	gFlags.errorLog = log.New(errorBuffer, "", 0)
	gFlags.exit = func(code int) {
		// THEN it should exit non-zero (failure)
		if code == 0 {
			t.Errorf("Unexpected exit code %v", code)
		}
	}

	// WHEN VerifyManifest is run
	main()

	// THEN it should fail with an error
	if errorBuffer.Len() == 0 {
		t.Error("An error should have been logged")
	}

	// THEN it shouldn't have got to output anything yet
	infoStr := infoBuffer.String()
	if strings.Contains(infoStr, "0cc175b9c0f1b6a831c399e269772661") {
		t.Error("Nothing should be output, because it should have failed.")
	}
}
