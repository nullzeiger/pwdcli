// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pwdhandling containing password handling functions
package pwdhandling

import (
	"fmt"
	"strings"

	"github.com/nullzeiger/pwdcli/internal/filehandling"
)

type Pwd = filehandling.Account

// All returns formatted entries
func All() ([]string, error) {
	accounts, err := filehandling.ReadJSON()
	if err != nil {
		return nil, err
	}

	entries := []string{}
	for i, acc := range accounts {
		entries = append(entries,
			fmt.Sprintf("[%d] Website: %s Username: %s Email: %s Password: %s",
				i, acc.Website, acc.Username, acc.Email, acc.Pwd))
	}
	return entries, nil
}

// Create adds a new entry
func Create(pwd Pwd) error {
	return filehandling.AppendJSON(pwd)
}

// Delete removes an entry by index
func Delete(index int) (bool, error) {
	accounts, err := filehandling.ReadJSON()
	if err != nil {
		return false, err
	}

	if index < 0 || index >= len(accounts) {
		return false, fmt.Errorf("index out of range")
	}

	accounts = append(accounts[:index], accounts[index+1:]...)
	return true, filehandling.WriteJSON(accounts)
}

// Search returns all accounts matching the keyword (case-insensitive)
// It returns a slice of structs containing the index and the account itself.
func Search(key string) ([]struct {
	Index   int
	Account Pwd
}, error,
) {
	accounts, err := filehandling.ReadJSON()
	if err != nil {
		return nil, err
	}

	key = strings.ToLower(key)
	results := []struct {
		Index   int
		Account Pwd
	}{}

	for i, acc := range accounts {
		if strings.Contains(strings.ToLower(acc.Website), key) ||
			strings.Contains(strings.ToLower(acc.Username), key) ||
			strings.Contains(strings.ToLower(acc.Email), key) ||
			strings.Contains(strings.ToLower(acc.Pwd), key) {

			results = append(results, struct {
				Index   int
				Account Pwd
			}{Index: i, Account: acc})
		}
	}

	return results, nil
}
