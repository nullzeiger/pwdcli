// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package handling provides higher-level business logic for managing
// password entries. It sits between the CLI layer and the storage layer,
// offering operations such as listing, creating, deleting, and searching
// account entries.
package handling

import (
	"fmt"
	"strings"

	"github.com/nullzeiger/pwdcli/internal/account"
	"github.com/nullzeiger/pwdcli/internal/storage"
)

// Act is an alias to account.Account for convenience within this package.
type Act = account.Account

// All retrieves all stored accounts and returns them formatted as strings,
// each containing index and field details. It is used primarily by the CLI
// when listing entries.
func All() ([]string, error) {
	accounts, err := storage.Read()
	if err != nil {
		return nil, err
	}

	entries := []string{}
	for i, acc := range accounts {
		// Format each account as a readable CLI entry.
		entries = append(entries,
			fmt.Sprintf("[%d] Website: %s Username: %s Email: %s Password: %s",
				i, acc.Website, acc.Username, acc.Email, acc.Pwd))
	}
	return entries, nil
}

// Create appends a new account entry to the storage file.
// It performs no validation—validation should be done at the CLI or higher layer.
func Create(act Act) error {
	return storage.Append(act)
}

// Delete removes an account by its index. It returns true if the operation
// succeeds, and an error if the index is invalid or storage access fails.
func Delete(index int) (bool, error) {
	accounts, err := storage.Read()
	if err != nil {
		return false, err
	}

	// Validate index bounds
	if index < 0 || index >= len(accounts) {
		return false, fmt.Errorf("index out of range")
	}

	// Remove the entry using slice manipulation.
	accounts = append(accounts[:index], accounts[index+1:]...)

	return true, storage.Write(accounts)
}

// Search scans all stored accounts and returns those matching the given
// keyword (case-insensitive). It compares the keyword with the website,
// username, email, and password fields.
//
// The result is a slice of structs containing both the index of the match
// and a copy of the corresponding account.
func Search(key string) ([]struct {
	Index   int
	Account Act
}, error) {

	accounts, err := storage.Read()
	if err != nil {
		return nil, err
	}

	key = strings.ToLower(key)

	results := []struct {
		Index   int
		Account Act
	}{}

	// Match accounts based on any field containing the keyword.
	for i, acc := range accounts {
		if strings.Contains(strings.ToLower(acc.Website), key) ||
			strings.Contains(strings.ToLower(acc.Username), key) ||
			strings.Contains(strings.ToLower(acc.Email), key) ||
			strings.Contains(strings.ToLower(acc.Pwd), key) {

			results = append(results, struct {
				Index   int
				Account Act
			}{Index: i, Account: acc})
		}
	}

	return results, nil
}
