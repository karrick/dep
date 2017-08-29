// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !windows

package gps

import (
	"os"
	"path/filepath"

	"github.com/karrick/godirwalk"
)

func stripVendor(path string, de *godirwalk.Dirent) error {
	if de.Name() == "vendor" {
		// ??? needless call to Lstat
		if _, err := os.Lstat(path); err != nil {
			return err
		}

		if de.IsSymlink() {
			referentInfo, err := os.Stat(path)
			if err != nil {
				return err
			}
			if referentInfo.IsDir() {
				return os.Remove(path)
			}
		}

		if de.IsDir() {
			if err := removeAll(path); err != nil {
				return err
			}
			return filepath.SkipDir
		}

	}

	return nil
}
