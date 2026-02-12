// Copyright 2026 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package storage provides low-level functions for creating, reading,
// writing, and appending account data to the JSON storage file.
// The file path and permissions are managed via the util package.
package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/nullzeiger/pwdcli/internal/account"
	"github.com/nullzeiger/pwdcli/internal/util"
)

// masterPassword holds the session password used to encrypt
// and decrypt the storage file. It must be set before any
// read or write operations.
var masterPassword string

// SetMasterPassword sets the password used for encrypting
// and decrypting the storage file.
func SetMasterPassword(pw string) {
	masterPassword = pw
}

// Create initializes the storage file if it does not already exist.
// It ensures that the file path returned by util.FilePath()
// exists and contains an empty encrypted JSON array ([]).
// If the file already exists, the function does nothing.
func Create() error {
	path := util.FilePath()

	// If the storage file already exists, nothing needs to be done.
	if util.FileExists(path) {
		return nil
	}

	if masterPassword == "" {
		return errors.New("master password not set")
	}

	// Create an empty JSON array and encrypt it.
	emptyData, err := json.Marshal([]account.Account{})
	if err != nil {
		return err
	}

	encrypted, err := encrypt(emptyData, masterPassword)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encrypted, util.Perm)
}

// Read loads all stored accounts from the encrypted storage file into a slice.
// Returns an error if the file cannot be read, decrypted, or the JSON is malformed.
func Read() ([]account.Account, error) {
	path := util.FilePath()

	if masterPassword == "" {
		return nil, errors.New("master password not set")
	}

	// Read raw file contents.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Decrypt file contents.
	decrypted, err := decrypt(data, masterPassword)
	if err != nil {
		return nil, err
	}

	// Decode JSON into a slice of Account structs.
	var accounts []account.Account
	err = json.Unmarshal(decrypted, &accounts)
	return accounts, err
}

// Write replaces the entire encrypted storage file with the provided slice of accounts.
// The file is overwritten using the permissions defined in util.Perm.
func Write(accounts []account.Account) error {
	path := util.FilePath()

	if masterPassword == "" {
		return errors.New("master password not set")
	}

	// Encode accounts as JSON.
	jsonData, err := json.Marshal(accounts)
	if err != nil {
		return err
	}

	// Encrypt the JSON data.
	encrypted, err := encrypt(jsonData, masterPassword)
	if err != nil {
		return err
	}

	// Overwrite the storage file with encrypted data.
	return os.WriteFile(path, encrypted, util.Perm)
}

// Append reads the existing accounts from storage, adds the new account,
// and writes all accounts back to the encrypted file.
// Returns an error if reading or writing fails.
func Append(acc account.Account) error {
	// Load existing data.
	accounts, err := Read()
	if err != nil {
		return err
	}

	// Add the new entry.
	accounts = append(accounts, acc)

	// Save updated account list.
	return Write(accounts)
}

// deriveKey generates a 32-byte AES key from the provided password and salt.
func deriveKey(password string, salt []byte) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	h.Write(salt)
	return h.Sum(nil)
}

// encrypt encrypts the given data using AES-GCM with a key derived from the password.
// The output format is: [16 bytes salt][nonce][ciphertext]
func encrypt(data []byte, password string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	result := append(salt, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// decrypt decrypts data produced by the encrypt function using the provided password.
func decrypt(data []byte, password string) ([]byte, error) {
	if len(data) < 16 {
		return nil, errors.New("invalid encrypted data")
	}

	salt := data[:16]
	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if len(data) < 16+nonceSize {
		return nil, errors.New("invalid encrypted data length")
	}

	nonce := data[16 : 16+nonceSize]
	ciphertext := data[16+nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}
