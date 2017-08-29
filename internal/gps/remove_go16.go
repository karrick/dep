// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.7

package gps

import (
	"os"
	"runtime"

	"github.com/karrick/godirwalk"
)

// removeAll removes path and any children it contains. It deals correctly with
// removal on Windows where, prior to Go 1.7, there were issues when files were
// set to read-only.
func removeAll(path string) error {
	// Only need special handling for windows
	if runtime.GOOS != "windows" {
		return os.RemoveAll(path)
	}

	// Simple case: if Remove works, we're done.
	err := os.Remove(path)
	if err == nil || os.IsNotExist(err) {
		return nil
	}

	// make sure all files are writable so we can delete them
	err = godirwalk.Walk(path, &godirwalk.Options{Callback: func(path string, _ *godirwalk.Dirent) error {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		mode := info.Mode()
		mode0200 := mode | 0200
		if mode0200 == mode {
			return nil // node is already writable
		}
		return os.Chmod(path, mode0200)
	}})
	if err != nil {
		return err
	}

	return os.Remove(path)
}
