// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cli provides command-line interface utilities for pwd management
package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/nullzeiger/pwdcli/internal/filehandling"
	"github.com/nullzeiger/pwdcli/internal/pwdhandling"
)

func Run() {
	// Commands
	listFlag := flag.Bool("all", false, "List all password entries")
	addFlag := flag.Bool("add", false, "Add a new password entry")
	deleteFlag := flag.Int("delete", -1, "Delete an entry by index")
	searchFlag := flag.String("search", "", "Search entries by keyword")

	// Fields for -add
	website := flag.String("website", "", "Website (required for -add)")
	username := flag.String("username", "", "Username (required for -add)")
	email := flag.String("email", "", "Email (required for -add)")
	password := flag.String("pwd", "", "Password (required for -add)")

	flag.Parse()

	// Ensure ~/.passwords.json exists
	if err := filehandling.Create(); err != nil {
		fmt.Println("Error creating password file:", err)
		return
	}

	// LIST
	if *listFlag {
		entries, err := pwdhandling.All()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, e := range entries {
			fmt.Println(e)
		}
		return
	}

	// ADD
	if *addFlag {
		if *website == "" || *username == "" || *email == "" || *password == "" {
			fmt.Println("Missing fields for -add: --website --username --email --pwd")
			os.Exit(1)
		}

		newEntry := pwdhandling.Pwd{
			Website:  *website,
			Username: *username,
			Email:    *email,
			Pwd:      *password,
		}

		if err := pwdhandling.Create(newEntry); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Entry added successfully.")
		return
	}

	// DELETE
	if *deleteFlag >= 0 {
		ok, err := pwdhandling.Delete(*deleteFlag)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if ok {
			fmt.Printf("Entry [%d] deleted.\n", *deleteFlag)
		}
		return
	}

	// SEARCH
	if *searchFlag != "" {
		matches, err := pwdhandling.Search(*searchFlag)
		if err != nil {
			fmt.Println("Error:", err)
		}

		if len(matches) == 0 {
			fmt.Println("No results found.")
			return
		}

		for _, m := range matches {
			fmt.Printf("[%d] Website: %s Username: %s Email: %s Password: %s\n",
				m.Index, m.Account.Website, m.Account.Username, m.Account.Email, m.Account.Pwd)
		}

	}

	flag.Usage()
}
