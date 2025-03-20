package main

import (
	"database-test/cmd/root"
	"database-test/cmd/seed"
)

func main() {
	// Add subcommands to the root command
	root.RootCmd.AddCommand(seed.Command)

	// Execute the root command
	root.Execute()
}
