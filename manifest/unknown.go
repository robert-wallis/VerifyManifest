// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved

package manifest

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// UnknownHashes is a list of `HashLocation`s
type UnknownHashes map[string]HashLocation

// HashLocation contains placeholder information for where a hash was located in a file.
// This information is used when showing an error so you can know from where in the manifest file it came.
type HashLocation struct {
	LineNumber int
	Line       string
}

// LoadUnknownHashes loads an unknown text file format, and look for strings that look like hashes.
func LoadUnknownHashes(filename string) (*UnknownHashes, error) {
	lineNumber := 1
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Couldn't scan file for hashes in %v: %v", filename, err)
	}
	defer file.Close()

	result := &UnknownHashes{}

	rd := bufio.NewReader(file)
	line := []rune{}
	for {
		r, _, err := rd.ReadRune()
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
			result.Set(lower, hl)
		}
		if sha1 := SHA1InString(line); sha1 != nil {
			lower := strings.ToLower(*sha1)
			hl := HashLocation{
				LineNumber: lineNumber,
				Line:       string(line),
			}
			result.Set(lower, hl)
		}
		line = []rune{}
		lineNumber++
	}
	return result, nil
}

// Set adds the hash to the list
func (u *UnknownHashes) Set(hash string, location HashLocation) {
	(*u)[hash] = location
}

// Get returns the hash from the list
func (u *UnknownHashes) Get(hash string) (location HashLocation, ok bool) {
	v, ok := (*u)[hash]
	return v, ok
}

// Remove removes the hash from the list
func (u *UnknownHashes) Remove(hash string) {
	delete(*u, hash)
}

// RemoveSum removes the hash from the list if it's in the Sum structure
func (u *UnknownHashes) RemoveSum(sum Sum) {
	if _, ok := u.Get(sum.MD5); ok {
		u.Remove(sum.MD5)
	}
	if _, ok := u.Get(sum.SHA1); ok {
		u.Remove(sum.SHA1)
	}
}

// MD5InString returns the substring that looks like an MD5 sum in a string
func MD5InString(line []rune) *string {
	return hexString(line, 32)
}

// SHA1InString returns the substring that looks like a SHA1 sum in a string
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
