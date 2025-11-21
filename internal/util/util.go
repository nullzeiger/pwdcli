// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package util provides utility functions and constants used across the application,
// including file path generation and file existence checks.
package util

import "os"

const (
	// Filename defines the default name of the JSON storage file.
	Filename = ".passwords.json"

	// Perm specifies the file permissions used when writing the storage file.
	// 0o644 = owner read/write, group read, others read.
	Perm = 0o644
)

// FilePath returns the full path of the storage file located
// in the user's home directory. If the home directory cannot be
// determined, the function panics since the application cannot continue.
func FilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Panic is acceptable here because without a home directory
		// the program cannot determine where to store user data.
		panic(err)
	}
	return home + "/" + Filename
}

// FileExists checks whether a file exists at the given path.
// It returns true only if os.Stat reports no error.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
