// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package account defines the Account data structure used to represent
// stored user credentials, including website, username, email, and password.
package account

// Account represents a single credential entry stored in the application.
// Each field includes a JSON tag to ensure proper encoding and decoding
// when saving or loading accounts from the storage file.
type Account struct {
	// Website is the domain or service the credentials belong to.
	Website string `json:"website"`

	// Username is the login username associated with the account.
	Username string `json:"username"`

	// Email is the email address linked to the account, if applicable.
	Email string `json:"email"`

	// Pwd stores the password for the account.
	Pwd string `json:"pwd"`
}
