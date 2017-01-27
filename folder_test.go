package main

import "testing"

func Test_streamHashes(t *testing.T) {
	// GIVEN the test folder and a stream of files
	dirname := "test_data"
	files := make(chan string)
	go func() {
		files <- "a.txt"
		files <- "b.txt"
		close(files)
	}()

	// WHEN the files are streamed
	results := make(chan *fileNameSum)
	go streamHashes(dirname, files, results)

	// THEN the hashes are generated for a
	a := <-results
	if a.Sum.MD5 != "0cc175b9c0f1b6a831c399e269772661" {
		t.Fatalf("a.txt hash unexpected %v", a)
	}

	// THEN the hashes are generated for b
	b := <-results
	if b.Sum.MD5 != "92eb5ffee6ae2fec3ad71c777531578f" {
		t.Fatalf("b.txt hash unexpected %v", b)
	}

	// THEN the channel is closed
	end := <-results
	if end != nil {
		t.Fatalf("Should have ended the stream %v", end)
	}
}

func Test_streamHashes_FileError(t *testing.T) {
	// GIVEN a folder that does exist, but a file that doesn't
	dirname := "test_data"
	files := make(chan string)
	go func() {
		files <- "noexist"
		files <- "b.txt"
		close(files)
	}()

	// WHEN the files are streamed
	results := make(chan *fileNameSum)
	err := streamHashes(dirname, files, results)

	// THEN an error should have occurred
	if err == nil {
		t.Error("Stream should have failed, because file didn't exist.")
	}

	for f := range results {
		t.Errorf("No results, but got %v", f)
	}
}
