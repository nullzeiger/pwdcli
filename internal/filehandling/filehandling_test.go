// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filehandling

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func withTempHome(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	err := os.Setenv("HOME", tmpDir)
	if err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}

	return tmpDir
}

func TestCreate(t *testing.T) {
	tmpHome := withTempHome(t)
	path := filepath.Join(tmpHome, filename)

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("test setup failure: file %s already exists", path)
	}

	if err := Create(); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to be created", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed reading created file: %v", err)
	}

	if string(data) != "[]" {
		t.Fatalf("expected file content '[]', got: %s", string(data))
	}
}

func TestWriteAndReadJSON(t *testing.T) {
	withTempHome(t) // home override only, no variable needed

	accounts := []Account{
		{Website: "example.com", Username: "john", Email: "john@example.com", Pwd: "123"},
		{Website: "test.com", Username: "alice", Email: "alice@test.com", Pwd: "456"},
	}

	if err := WriteJSON(accounts); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	readAccounts, err := ReadJSON()
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	if len(readAccounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(readAccounts))
	}

	if readAccounts[0].Website != "example.com" {
		t.Fatalf("unexpected first item: %+v", readAccounts[0])
	}
}

func TestAppendJSON(t *testing.T) {
	withTempHome(t)

	if err := Create(); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	account := Account{
		Website:  "github.com",
		Username: "dev",
		Email:    "dev@example.com",
		Pwd:      "pass",
	}

	if err := AppendJSON(account); err != nil {
		t.Fatalf("AppendJSON failed: %v", err)
	}

	accounts, err := ReadJSON()
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	if len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accounts))
	}

	if accounts[0].Email != "dev@example.com" {
		t.Fatalf("unexpected appended account: %+v", accounts[0])
	}
}

func TestFileExists(t *testing.T) {
	tmpHome := withTempHome(t)
	path := filepath.Join(tmpHome, filename)

	if fileExists(path) {
		t.Fatalf("expected false for non-existing file")
	}

	if err := os.WriteFile(path, []byte("[]"), 0o644); err != nil {
		t.Fatalf("test setup failed: %v", err)
	}

	if !fileExists(path) {
		t.Fatalf("expected true for existing file")
	}
}

func TestJSONStructure(t *testing.T) {
	tmpHome := withTempHome(t)

	acc := []Account{{Website: "a", Username: "b", Email: "c", Pwd: "d"}}
	if err := WriteJSON(acc); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(tmpHome, filename))

	var decoded []map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("JSON unmarshalling failed: %v", err)
	}

	if _, ok := decoded[0]["website"]; !ok {
		t.Fatalf("missing JSON key 'website'")
	}
	if _, ok := decoded[0]["username"]; !ok {
		t.Fatalf("missing JSON key 'username'")
	}
}
