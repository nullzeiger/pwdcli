// Copyright 2026 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cli provides the command-line interface for interacting
// with the password management system.
//
// It handles flag parsing, user input/output, and coordinates
// data flow between the handling logic and the storage layers.
package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	handling "github.com/nullzeiger/pwdcli/internal/handling"
	"github.com/nullzeiger/pwdcli/internal/storage"
)

// Run is the main entry point for the CLI application.
// It defines command-line flags, initializes the storage environment,
// and dispatches actions based on the provided arguments.
func Run() {
	// Flag Definitions.
	listFlag := flag.Bool("all", false, "List all password entries")
	addFlag := flag.Bool("add", false, "Add a new password entry")
	deleteFlag := flag.Int("delete", -1, "Delete an entry by index")
	searchFlag := flag.String("search", "", "Search entries by keyword")

	// Fields required for adding a new entry.
	website := flag.String("website", "", "Website URL or name (required for -add)")
	username := flag.String("username", "", "Account username (required for -add)")
	email := flag.String("email", "", "Account email (required for -add)")
	password := flag.String("pwd", "", "Account password (required for -add)")

	flag.Parse()

	// If no operational flag is provided, display usage information and exit.
	if !*listFlag && !*addFlag && *deleteFlag < 0 && *searchFlag == "" {
		flag.Usage()
		return
	}

	// Prompt for the master password to initialize storage access.
	fmt.Print("Master password: ")
	master, err := readPassword()
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	if master == "" {
		fmt.Println("Master password cannot be empty")
		return
	}

	storage.SetMasterPassword(master)

	// Create or verify the existence of the underlying storage file.
	if err := storage.Create(); err != nil {
		fmt.Println("Error creating password file:", err)
		return
	}

	// Execute LIST command.
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

	// Execute ADD command.
	if *addFlag {
		if *website == "" || *username == "" || *email == "" || *password == "" {
			fmt.Println("Missing fields for -add: --website --username --email --pwd")
			return
		}

		newEntry := handling.Act{
			Website:  *website,
			Username: *username,
			Email:    *email,
			Pwd:      *password,
		}

		if err := handling.Create(newEntry); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Entry added successfully.")
		return
	}

	// Execute DELETE command.
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

	// Execute SEARCH command.
	if *searchFlag != "" {
		matches, err := handling.Search(*searchFlag)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(matches) == 0 {
			fmt.Println("No results found.")
			return
		}

		for _, m := range matches {
			fmt.Printf(
				"[%d] Website: %s Username: %s Email: %s Password: %s\n",
				m.Index, m.Account.Website, m.Account.Username, m.Account.Email, m.Account.Pwd,
			)
		}
		return
	}
}

// readPassword reads a password from stdin without echoing characters to the terminal.
// It works by temporarily disabling the ECHO flag in the terminal settings.
// This function is Unix-specific (Linux, macOS).
// It manipulates terminal settings directly via syscalls and restores the original
// state after reading the input to ensure the terminal remains functional.
func readPassword() (string, error) {
	// Get the file descriptor for standard input.
	fd := int(os.Stdin.Fd())

	// Retrieve the current terminal state using the TCGETS ioctl command.
	// This captures the original terminal configuration so we can restore it later.
	var oldState syscall.Termios
	if _, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		syscall.TCGETS,
		uintptr(unsafe.Pointer(&oldState)),
		0, 0, 0,
	); errno != 0 {
		return "", fmt.Errorf("unable to get terminal state: %v", errno)
	}

	// Create a copy of the terminal state and disable the ECHO flag.
	// The ECHO flag controls whether input characters are displayed on the terminal.
	newState := oldState
	newState.Lflag &^= syscall.ECHO

	// Apply the modified terminal settings using the TCSETS ioctl command.
	// This immediately disables character echoing for the current terminal session.
	if _, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		syscall.TCSETS,
		uintptr(unsafe.Pointer(&newState)),
		0, 0, 0,
	); errno != 0 {
		return "", fmt.Errorf("unable to set terminal state: %v", errno)
	}

	// Ensure the original terminal state is restored after reading the password,
	// regardless of whether an error occurs. This prevents the terminal from
	// remaining in a non-echo state, which would make it unusable.
	defer func() {
		syscall.Syscall6(
			syscall.SYS_IOCTL,
			uintptr(fd),
			syscall.TCSETS,
			uintptr(unsafe.Pointer(&oldState)),
			0, 0, 0,
		)
		// Print a newline since the user's Enter key press was not echoed.
		fmt.Println()
	}()

	// Read the password input from the user.
	// The input will not be displayed on screen due to the disabled ECHO flag.
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Remove the trailing newline character(s) from the input.
	if len(password) > 0 && password[len(password)-1] == '\n' {
		password = password[:len(password)-1]
	}

	return password, nil
}
