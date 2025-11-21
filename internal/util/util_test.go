// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package util_test contains unit tests for the util package.
// These tests validate file path generation and file existence checks
// without affecting the user's real HOME directory or filesystem.
package util_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nullzeiger/pwdcli/internal/util"
)

// TestFilePath verifies that util.FilePath correctly constructs the
// full path to the password storage file based on the user's home directory.
// The test uses a temporary directory to safely override HOME.
func TestFilePath(t *testing.T) {
	// Create a temporary directory to simulate HOME
	tmpDir := t.TempDir()

	// Override HOME environment variable for the test
	t.Setenv("HOME", tmpDir)

	// Construct the expected file path
	expected := filepath.Join(tmpDir, util.Filename)

	// Call FilePath
	got := util.FilePath()

	// Compare result
	if got != expected {
		t.Fatalf("FilePath() = %s; want %s", got, expected)
	}
}

// TestFileExists checks that util.FileExists correctly identifies
// whether a file exists at a given path. It tests both existing and
// non-existing files using a temporary directory.
func TestFileExists(t *testing.T) {
	// Create temporary directory for the test
	tmpDir := t.TempDir()

	// Create a temporary file to test existence
	filePath := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// FileExists should return true for the existing file
	if !util.FileExists(filePath) {
		t.Fatalf("FileExists(%s) = false; want true", filePath)
	}

	// FileExists should return false for a non-existing file
	nonExistent := filepath.Join(tmpDir, "does_not_exist.txt")
	if util.FileExists(nonExistent) {
		t.Fatalf("FileExists(%s) = true; want false", nonExistent)
	}
}
