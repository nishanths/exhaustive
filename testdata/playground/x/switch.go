package x

import (
	"io/fs"
	"os"
)

func _q() {
	// Type alias:
	// type os.FileMode = fs.FileMode

	fi, _ := os.Lstat(".")

	switch fi.Mode() { // want "^missing cases in switch of type fs.FileMode: ModeDevice, ModePerm, ModeSetgid, ModeSetuid, ModeType$"
	case os.ModeDir: // should not report ModeDir, because os.ModeDir has the same name and same value as fs.ModeDir
	case os.ModeAppend: // "
	case os.ModeExclusive: // "
	case fs.ModeTemporary:
	case fs.ModeSymlink:
	case fs.ModeNamedPipe:
	case os.ModeSocket: // "
	case fs.ModeCharDevice:
	case fs.ModeSticky:
	case fs.ModeIrregular:
	}
}
