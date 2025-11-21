// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cli provides the command-line interface used for interacting
// with the password management system. It handles parsing flags,
// dispatching commands, and coordinating with the underlying handling
// and storage layers.
package cli

import (
	"flag"
	"fmt"
	"os"

	handling "github.com/nullzeiger/pwdcli/internal/handling"
	"github.com/nullzeiger/pwdcli/internal/storage"
)

// Run is the main entry point for the CLI. It defines and parses flags,
// ensures the storage file exists, and dispatches the appropriate action
// based on the user’s command-line arguments.
func Run() {
	// --- Command Flags ---
	// Basic operations
	listFlag := flag.Bool("all", false, "List all password entries")
	addFlag := flag.Bool("add", false, "Add a new password entry")
	deleteFlag := flag.Int("delete", -1, "Delete an entry by index")
	searchFlag := flag.String("search", "", "Search entries by keyword")

	// Fields required when using -add
	website := flag.String("website", "", "Website (required for -add)")
	username := flag.String("username", "", "Username (required for -add)")
	email := flag.String("email", "", "Email (required for -add)")
	password := flag.String("pwd", "", "Password (required for -add)")

	flag.Parse()

	// Ensure the storage file exists (~/.passwords.json)
	// If it doesn't, it is automatically created.
	if err := storage.Create(); err != nil {
		fmt.Println("Error creating password file:", err)
		return
	}

	// --- LIST COMMAND ---
	if *listFlag {
		entries, err := handling.All()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, e := range entries {
			fmt.Println(e)
		}
		return
	}

	// --- ADD COMMAND ---
	if *addFlag {
		// Validate required fields
		if *website == "" || *username == "" || *email == "" || *password == "" {
			fmt.Println("Missing fields for -add: --website --username --email --pwd")
			os.Exit(1)
		}

		// Construct new entry
		newEntry := handling.Act{
			Website:  *website,
			Username: *username,
			Email:    *email,
			Pwd:      *password,
		}

		// Save the new entry
		if err := handling.Create(newEntry); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Entry added successfully.")
		return
	}

	// --- DELETE COMMAND ---
	if *deleteFlag >= 0 {
		ok, err := handling.Delete(*deleteFlag)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if ok {
			fmt.Printf("Entry [%d] deleted.\n", *deleteFlag)
		}
		return
	}

	// --- SEARCH COMMAND ---
	if *searchFlag != "" {
		matches, err := handling.Search(*searchFlag)
		if err != nil {
			fmt.Println("Error:", err)
		}

		// No results found
		if len(matches) == 0 {
			fmt.Println("No results found.")
			return
		}

		// Print matching entries
		for _, m := range matches {
			fmt.Printf(
				"[%d] Website: %s Username: %s Email: %s Password: %s\n",
				m.Index, m.Account.Website, m.Account.Username, m.Account.Email, m.Account.Pwd,
			)
		}
		return
	}

	// If no command was matched, print usage help.
	flag.Usage()
}
