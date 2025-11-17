// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package filehandling containing file handling functions
package filehandling

import (
	"encoding/json"
	"os"
)

const (
	perm     = 0o644
	filename = ".passwords.json"
)

// Account represents the JSON structure stored in the file
type Account struct {
	Website  string `json:"website"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Pwd      string `json:"pwd"`
}

// getFilePath returns the absolute path to ~/.passwords.json
func getFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home + "/" + filename
}

// Create ensures ~/.passwords.json exists
func Create() error {
	path := getFilePath()

	if fileExists(path) {
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte("[]"))
	return err
}

// ReadJSON reads all accounts from ~/.passwords.json
func ReadJSON() ([]Account, error) {
	path := getFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var accounts []Account
	err = json.Unmarshal(data, &accounts)
	return accounts, err
}

// WriteJSON overwrites ~/.passwords.json
func WriteJSON(accounts []Account) error {
	path := getFilePath()

	jsonData, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, perm)
}

// AppendJSON adds a new account
func AppendJSON(account Account) error {
	accounts, err := ReadJSON()
	if err != nil {
		return err
	}

	accounts = append(accounts, account)
	return WriteJSON(accounts)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
