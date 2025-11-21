// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package handling_test contains unit tests for the handling package.
// These tests verify the correct behavior of high-level account operations
// such as creating, listing, deleting, and searching accounts.
package handling_test

import (
	"testing"

	"github.com/nullzeiger/pwdcli/internal/handling"
	"github.com/nullzeiger/pwdcli/internal/storage"
	"github.com/nullzeiger/pwdcli/internal/util"
)

// setupTempStorage prepares a temporary environment for tests.
// It overrides HOME to a temp directory and ensures the storage file exists.
// Returns the expected storage file path.
func setupTempStorage(t *testing.T) string {
	t.Helper() // Marks this as a helper; errors point to the caller
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Initialize the storage file
	if err := storage.Create(); err != nil {
		t.Fatalf("storage.Create() failed: %v", err)
	}

	return util.FilePath()
}

// TestCreateAndAll verifies that creating an account via handling.Create
// works correctly and that handling.All returns the correct formatted output.
func TestCreateAndAll(t *testing.T) {
	setupTempStorage(t)

	acc := handling.Act{Website: "example.com", Username: "user", Email: "a@b.com", Pwd: "123"}

	// Create a new account
	if err := handling.Create(acc); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Retrieve all entries
	entries, err := handling.All()
	if err != nil {
		t.Fatalf("All() failed: %v", err)
	}

	// Expect exactly one entry
	if len(entries) != 1 {
		t.Fatalf("All() returned %d entries; want 1", len(entries))
	}

	// Verify the formatted string
	expectedSubstring := "Website: example.com Username: user Email: a@b.com Password: 123"
	if entries[0] != "[0] "+expectedSubstring {
		t.Fatalf("All()[0] = %s; want %s", entries[0], "[0] "+expectedSubstring)
	}
}

// TestDelete verifies that handling.Delete removes accounts correctly
// and handles invalid indices properly.
func TestDelete(t *testing.T) {
	setupTempStorage(t)

	acc1 := handling.Act{Website: "site1", Username: "u1", Email: "e1", Pwd: "p1"}
	acc2 := handling.Act{Website: "site2", Username: "u2", Email: "e2", Pwd: "p2"}

	handling.Create(acc1)
	handling.Create(acc2)

	// Delete the second account
	ok, err := handling.Delete(1)
	if err != nil || !ok {
		t.Fatalf("Delete(1) failed: %v", err)
	}

	// Verify remaining accounts
	accounts, _ := storage.Read()
	if len(accounts) != 1 || accounts[0].Website != "site1" {
		t.Fatalf("After delete, remaining accounts = %v; want only site1", accounts)
	}

	// Attempt deletion with an invalid index
	ok, err = handling.Delete(10)
	if err == nil || ok {
		t.Fatalf("Delete(10) should fail for invalid index")
	}
}

// TestSearch verifies that handling.Search correctly finds accounts
// based on a keyword and performs case-insensitive matching.
func TestSearch(t *testing.T) {
	setupTempStorage(t)

	acc1 := handling.Act{Website: "google.com", Username: "user1", Email: "a@b.com", Pwd: "pass1"}
	acc2 := handling.Act{Website: "example.com", Username: "user2", Email: "c@d.com", Pwd: "pass2"}

	handling.Create(acc1)
	handling.Create(acc2)

	// Search by website keyword
	results, err := handling.Search("google")
	if err != nil {
		t.Fatalf("Search() failed: %v", err)
	}
	if len(results) != 1 || results[0].Account.Website != "google.com" {
		t.Fatalf("Search by 'google' returned wrong result: %v", results)
	}

	// Case-insensitive search
	results, _ = handling.Search("EXAMPLE")
	if len(results) != 1 || results[0].Account.Website != "example.com" {
		t.Fatalf("Case-insensitive search failed: %v", results)
	}

	// Search for non-existing keyword should return 0 results
	results, _ = handling.Search("notfound")
	if len(results) != 0 {
		t.Fatalf("Search for 'notfound' should return 0 results, got %d", len(results))
	}
}
