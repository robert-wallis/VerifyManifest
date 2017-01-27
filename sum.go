// Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
package main

import (
	"io"
	"fmt"
	"strings"
	"crypto/sha1"
	"crypto/md5"
	"os"
)

type Sum struct {
	MD5 string
	SHA1 string
}

func (s *Sum) Calculate(fileName string) error {
	sha1hash := sha1.New()
	md5hash := md5.New()
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buffer := make([]byte, 65536)
	for {
		count, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		sha1hash.Write(buffer[:count])
		md5hash.Write(buffer[:count])
	}
	s.SHA1 = fmt.Sprintf("%x", sha1hash.Sum(nil))
	s.MD5 = fmt.Sprintf("%x", md5hash.Sum(nil))
	return nil
}

func (s *Sum) Verify(other Sum) error {
	if strings.ToLower(s.MD5) != strings.ToLower(other.MD5) {
		return fmt.Errorf("MD5 mismatch %v != %v", s.MD5, other.MD5)
	}
	if strings.ToLower(s.SHA1) != strings.ToLower(other.SHA1) {
		return fmt.Errorf("SHA1 mismatch %v != %v", s.SHA1, other.SHA1)
	}
	return nil
}
