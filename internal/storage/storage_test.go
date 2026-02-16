// Copyright 2026 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package storage_test

import (
	"testing"

	"github.com/nullzeiger/pwdcli/internal/account"
	"github.com/nullzeiger/pwdcli/internal/storage"
	"github.com/nullzeiger/pwdcli/internal/util"
)

// setupTempStorage prepares a temporary storage environment for tests.
func setupTempStorage(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Set a master password for tests
	storage.SetMasterPassword([]byte("Passwordsupersicura"))

	path := util.FilePath()

	if err := storage.Create(); err != nil {
		t.Fatalf("storage.Create() failed: %v", err)
	}

	return path
}

// TestCreate verifies that storage.Create correctly creates the storage file.
func TestCreate(t *testing.T) {
	path := setupTempStorage(t)

	if !util.FileExists(path) {
		t.Fatalf("File %s should exist after Create()", path)
	}

	// We cannot check for "[]" anymore because the file is encrypted.
	// Instead we just verify that it can be read and decrypted correctly.

	accounts, err := storage.Read()
	if err != nil {
		t.Fatalf("Read after Create() failed: %v", err)
	}

	if len(accounts) != 0 {
		t.Fatalf("Expected empty account list, got %d entries", len(accounts))
	}
}

// TestWriteAndRead verifies that storage.Write correctly saves a slice of accounts
// and that storage.Read can read them back accurately.
func TestWriteAndRead(t *testing.T) {
	setupTempStorage(t)

	accounts := []account.Account{
		{Website: "example.com", Username: "user", Email: "a@b.com", Pwd: "123"},
	}

	if err := storage.Write(accounts); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	readAccounts, err := storage.Read()
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}

	if len(readAccounts) != 1 {
		t.Fatalf("Read returned %d accounts; want 1", len(readAccounts))
	}

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

	if err := storage.Append(acc1); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}

	if err := storage.Append(acc2); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}

	accounts, err := storage.Read()
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}

	if len(accounts) != 2 {
		t.Fatalf("Read returned %d accounts; want 2", len(accounts))
	}

	if accounts[0].Website != "example1.com" || accounts[1].Website != "example2.com" {
		t.Fatalf("Accounts data mismatch: %v", accounts)
	}
}

// TestAppendWithoutCreate verifies that Append fails if the storage file
// has not been created.
func TestAppendWithoutCreate(t *testing.T) {
	t.Helper()

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Even in this test we must set a master password
	storage.SetMasterPassword([]byte("Passwordsupersicura"))

	acc := account.Account{Website: "site.com", Username: "u", Email: "e", Pwd: "p"}

	err := storage.Append(acc)
	if err == nil {
		t.Fatalf("Append() should fail if storage file does not exist")
	}
}
