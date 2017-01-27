// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type HashLocation struct {
	LineNumber int
	Line       string
}

func LoadUnknownForHashes(filename string) (map[string]HashLocation, error) {
	lineNumber := 0
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Couldn't scan file for hashes %v: %v", filename, err)
	}
	defer file.Close()

	result := make(map[string]HashLocation)

	rd := bufio.NewReader(file)
	line := []rune{}
	i := 0
	for {
		r, c, err := rd.ReadRune()
		i += c
		if err != nil && err != io.EOF {
			err = fmt.Errorf("Error scanning for hashes in %v: %v", filename, err)
			return nil, err
		}
		if err == io.EOF {
			break
		}
		if r != '\n' {
			if r != 0 { // utf16 hack
				line = append(line, r)
			}
			continue
		}
		if md5 := MD5InString(line); md5 != nil {
			lower := strings.ToLower(*md5)
			hl := HashLocation{
				LineNumber: lineNumber,
				Line:       string(line),
			}
			result[lower] = hl
		}
		if sha1 := SHA1InString(line); sha1 != nil {
			lower := strings.ToLower(*sha1)
			hl := HashLocation{
				LineNumber: lineNumber,
				Line:       string(line),
			}
			result[lower] = hl
		}
		line = []rune{}
		lineNumber++
	}
	return result, nil
}

// returns the substring that looks like an MD5 sum in a string
func MD5InString(line []rune) *string {
	return hexString(line, 32)
}

// returns the substring that looks like a SHA1 sum in a string
func SHA1InString(line []rune) *string {
	return hexString(line, 40)
}

// returns true if the rune is a valid hexadecimal character
func hexRune(r rune) bool {
	if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
		return true
	}
	return false
}

// search a string for a specific length of hex characters
func hexString(line []rune, length int) *string {
	start := 0
	count := 0
	result := ""
	for l := range line {
		ch := line[l]
		if hexRune(ch) {
			count++
			if count == length {
				if l == len(line)-1 || !hexRune(line[l+1]) {
					result = string(line[start : start+length])
					return &result
				}
			}
		} else {
			start = l + 1
			count = 0
		}
	}
	return nil
}
