package x

import (
	"io/fs"
	"os"
)

func _q() {
	// Type alias:
	// type os.FileMode = fs.FileMode
	//
	// Of interest, note that e.g. listing os.ModeSocket in a case clause is
	// equivalent to listing fs.ModeSocket (both have the same constant value).

	var fi fs.FileInfo

	switch fi.Mode() { // want "^missing cases in switch of type fs.FileMode: fs.ModeDevice, fs.ModeSetuid, fs.ModeSetgid, fs.ModeType, fs.ModePerm$"
	case os.ModeDir:
	case os.ModeAppend:
	case os.ModeExclusive:
	case fs.ModeTemporary:
	case fs.ModeSymlink:
	case fs.ModeNamedPipe, os.ModeSocket:
	case fs.ModeCharDevice:
	case fs.ModeSticky:
	case fs.ModeIrregular:
	}

	_ = map[fs.FileMode]int{ // want "^missing keys in map of key type fs.FileMode: fs.ModeDevice, fs.ModeSetuid, fs.ModeSetgid, fs.ModeType, fs.ModePerm$"
		os.ModeDir:        1,
		os.ModeAppend:     2,
		os.ModeExclusive:  3,
		fs.ModeTemporary:  4,
		fs.ModeSymlink:    5,
		fs.ModeNamedPipe:  6,
		os.ModeSocket:     7,
		fs.ModeCharDevice: 8,
		fs.ModeSticky:     9,
		fs.ModeIrregular:  10,
	}
}
