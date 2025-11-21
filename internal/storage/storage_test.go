// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package storage_test contains unit tests for the storage package.
// These tests validate creating, reading, writing, and appending account
// data using temporary directories to avoid affecting the user's real data.
package storage_test

import (
	"os"
	"testing"

	"github.com/nullzeiger/pwdcli/internal/account"
	"github.com/nullzeiger/pwdcli/internal/storage"
	"github.com/nullzeiger/pwdcli/internal/util"
)

// setupTempStorage prepares a temporary storage environment for tests.
// It performs the following steps:
// - Creates a temporary directory to simulate a HOME directory
// - Overrides the HOME environment variable for the test
// - Ensures that the storage file exists and is initialized as empty
func setupTempStorage(t *testing.T) string {
	t.Helper() // Marks this function as a helper to report errors at the caller

	// Create a temporary directory for this test
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Get the full path to the storage file
	path := util.FilePath()

	// Create the storage file if it does not exist
	if err := storage.Create(); err != nil {
		t.Fatalf("storage.Create() failed: %v", err)
	}

	return path
}

// TestCreate verifies that storage.Create correctly creates the storage file
// and initializes it with an empty JSON array ("[]").
func TestCreate(t *testing.T) {
	path := setupTempStorage(t)

	// Confirm that the file now exists
	if !util.FileExists(path) {
		t.Fatalf("File %s should exist after Create()", path)
	}

	// Verify the content of the file
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(data) != "[]" {
		t.Fatalf("File content = %s; want []", string(data))
	}
}

// TestWriteAndRead verifies that storage.Write correctly saves a slice of accounts
// and that storage.Read can read them back accurately.
func TestWriteAndRead(t *testing.T) {
	setupTempStorage(t)

	accounts := []account.Account{
		{Website: "example.com", Username: "user", Email: "a@b.com", Pwd: "123"},
	}

	// Write the account to storage
	if err := storage.Write(accounts); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	// Read accounts back from storage
	readAccounts, err := storage.Read()
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}

	// Verify that exactly one account was read
	if len(readAccounts) != 1 {
		t.Fatalf("Read returned %d accounts; want 1", len(readAccounts))
	}

	// Check the account data matches what was written
	if readAccounts[0].Website != "example.com" {
		t.Fatalf("Account Website = %s; want example.com", readAccounts[0].Website)
	}
}

// TestAppend verifies that storage.Append correctly adds new accounts
// without overwriting previous entries.
func TestAppend(t *testing.T) {
	setupTempStorage(t)

	acc1 := account.Account{Website: "example1.com", Username: "u1", Email: "e1", Pwd: "p1"}
	acc2 := account.Account{Website: "example2.com", Username: "u2", Email: "e2", Pwd: "p2"}

	// Append the first account
	if err := storage.Append(acc1); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}

	// Append the second account
	if err := storage.Append(acc2); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}

	// Read all accounts from storage
	accounts, err := storage.Read()
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}

	// Check that both accounts exist
	if len(accounts) != 2 {
		t.Fatalf("Read returned %d accounts; want 2", len(accounts))
	}

	// Verify that the accounts are in the correct order
	if accounts[0].Website != "example1.com" || accounts[1].Website != "example2.com" {
		t.Fatalf("Accounts data mismatch: %v", accounts)
	}
}

// TestAppendWithoutCreate verifies that Append fails if the storage file
// has not been created, simulating incorrect usage.
func TestAppendWithoutCreate(t *testing.T) {
	t.Helper()

	// Create temporary HOME directory but do not create storage file
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	acc := account.Account{Website: "site.com", Username: "u", Email: "e", Pwd: "p"}

	// Append should fail because the storage file does not exist
	err := storage.Append(acc)
	if err == nil {
		t.Fatalf("Append() should fail if storage file does not exist")
	}
}
