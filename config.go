package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetConfigDir() string {
	if runtime.GOOS == "windows" || os.Geteuid() != 0 {
		// running on windows (-1) or not as root → return process dir
		e, err := os.Executable()
		if err != nil {
			return "." // no better option
		}
		return filepath.Dir(e)
	}

	p := filepath.FromSlash("/etc/seidan")
	if i, err := os.Stat(p); os.IsNotExist(err) {
		// file does not exist
		err = os.Mkdir(p, 0700)
		if err != nil {
			// this is bad
			panic(err)
		}
	} else if err != nil {
		// failed to stat, but not file not exist → this is bad
		panic(err)
	} else if !i.IsDir() {
		panic("/etc/seidan exists and is not a directory")
	}

	return p
}
