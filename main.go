// Copyright 2025 Ivan Guerreschi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package main is the entry point of the password manager application.
// It delegates execution to the CLI layer, which handles command parsing
// and interactions with the underlying storage and handling logic.
package main

import "github.com/nullzeiger/pwdcli/internal/cli"

// main initializes the command-line interface by calling cli.Run(),
// which processes flags and executes the corresponding action.
func main() {
	cli.Run()
}
