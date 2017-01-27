# Verify Manifest

Helps create a new `manifest.json` file in a folder, and verify an existing manifest.
Each file is checked for `SHA1` and `MD5` hashes.

This is helpful to see if the contents of a file have changed since the last time the tool was run.

# Installation

Download to your `PATH`
* [VerifyManifest.exe v0.1](https://github.com/robert-wallis/VerifyManifest/releases/download/v0.1/VerifyManifest.exe)

### Install from Source
```
go install github.com/robert-wallis/VerifyManifest
```

# Usage
```
D:> VerifyManifest.exe
a.txt	md5:0cc175b9c0f1b6a831c399e269772661	sha1:86f7e437faa5a7fce15d1ddcb9eaeaea377667b8
b.txt	md5:92eb5ffee6ae2fec3ad71c777531578f	sha1:e9d71f5ee7c92d6dc9e92ffdad17b8bd49418f98
Saved manifest to manifest.json
```

```manifest.json
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