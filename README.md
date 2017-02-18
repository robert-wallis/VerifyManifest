[![goreportcard.com](https://goreportcard.com/badge/github.com/robert-wallis/VerifyManifest)](https://goreportcard.com/report/github.com/robert-wallis/VerifyManifest)

# Verify Manifest

Helps create and verify the SHA1 and MD5 hashes :1234: with a `manifest.json` file.
If the contents of a folder change, then VerifyManifest will tell you when you run it.

# Installation

### Download
Download [VerifyManifest.exe v0.2](https://github.com/robert-wallis/VerifyManifest/releases/download/v0.2/VerifyManifest.exe) to your `PATH`somewhere.

### Or
### Install from Source
```
go install github.com/robert-wallis/VerifyManifest
```
`go install` automatically puts the VerifyManifest.exe in your `$GOPATH/bin` and you should have your `$GOPATH/bin` in your `$PATH` so you can install other go based tools.

# Usage
When you run `VerifyManifest` it will calculate hashes for all the files in the folder and save them to `manifest.json`
```
D:\test_data> VerifyManifest
a.txt	md5:0cc175b9c0f1b6a831c399e269772661	sha1:86f7e437faa5a7fce15d1ddcb9eaeaea377667b8
b.txt	md5:92eb5ffee6ae2fec3ad71c777531578f	sha1:e9d71f5ee7c92d6dc9e92ffdad17b8bd49418f98
Saved manifest to manifest.json
```

`manifest.json`:
```json
{
	"a.txt": {
		"MD5": "0cc175b9c0f1b6a831c399e269772661",
		"SHA1": "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8"
	},
	"b.txt": {
		"MD5": "92eb5ffee6ae2fec3ad71c777531578f",
		"SHA1": "e9d71f5ee7c92d6dc9e92ffdad17b8bd49418f98"
	}
}
```

Now when the data changes, and VerifyManifest is run again, it will report an error.  It will not save to `manifest.json` unless all the existing hashes are successfully verified.
```
D:\test_data> echo z > b.txt
D:\test_data> VerifyManifest
a.txt   md5:0cc175b9c0f1b6a831c399e269772661    sha1:86f7e437faa5a7fce15d1ddcb9eaeaea377667b8
Error b.txt: MD5 mismatch 92eb5ffee6ae2fec3ad71c777531578f != efaddc0ff690c7f1f7d802143b5172be
b.txt   md5:efaddc0ff690c7f1f7d802143b5172be    sha1:b234c9cbc82c27e7f996dd4744791336ed5ea287
```

#### Copyright (C) 2017 Robert A. Wallis, All Rights Reserved
