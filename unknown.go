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

func LoadUnknownForHashes(filename string) (result map[string]HashLocation, err error) {
	//lineNumber := 0
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Couldn't scan file for hashes %v: %v", filename, err)
	}

	rd := bufio.NewReader(file)
	for {
		_, err = rd.ReadString('\n')
		if err != nil && err != io.EOF {
			err = fmt.Errorf("Error scanning for hashes in %v: %v", filename, err)
			return
		}
		if err == io.EOF {
			break
		}
	}
	return
}

// returns the substring that looks like an MD5 sum in a string
func MD5InString(line string) *string {
	return hexString(line, 32)
}

// returns the substring that looks like a SHA1 sum in a string
func SHA1InString(line string) *string {
	return hexString(line, 40)
}

// returns true if the rune is a valid hexidecimal character
func hexRune(r rune) bool {
	if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
		return true
	}
	return false
}

// search a string for a specific length of hex characters
func hexString(line string, length int) *string {
	start := 0
	count := 0
	i := 0
	rd := strings.NewReader(line)
	for {
		ch, _, err := rd.ReadRune()
		if err != nil {
			break
		}
		if hexRune(ch) {
			count++
			if count == length {
				ch, _, err := rd.ReadRune()
				if err != nil || !hexRune(ch) {
					result := line[start: start + length]
					return &result
				}
				rd.UnreadRune()
			}
		} else {
			start = i + 1
			count = 0
		}
		i++
	}
	return nil
}
