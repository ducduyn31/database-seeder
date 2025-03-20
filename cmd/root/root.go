package root

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

var (
	// Used for flags
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
	dbSSLMode  string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dbseeder",
	Short: "A CLI tool for seeding e-commerce database with fake data",
	Long: `Database Seeder is a CLI application that generates and inserts
fake data into your e-commerce database. It can create users, products,
categories, orders, and more using realistic fake data.

This tool is designed to help developers quickly populate their
database with test data for development and testing purposes.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Check if we're running the completion command or __complete (used by shell completion)
	isCompletion := false

	for _, arg := range os.Args {
		if arg == "completion" || arg == "__complete" {
			isCompletion = true
			break
		}
	}

	// Show the banner if not running completion
	if !isCompletion {
		// Print a fancy banner
		pterm.DefaultBigText.WithLetters(
			putils.LettersFromStringWithStyle("DB", pterm.NewStyle(pterm.FgLightCyan)),
			putils.LettersFromStringWithStyle("Seeder", pterm.NewStyle(pterm.FgLightGreen)),
		).Render()

		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgDarkGray)).Println("E-Commerce Database Seeder")
		pterm.Println()
	}

	if err := RootCmd.Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&dbHost, "host", "localhost", "Database host")
	RootCmd.PersistentFlags().IntVar(&dbPort, "port", 5433, "Database port")
	RootCmd.PersistentFlags().StringVar(&dbUser, "user", "shared_user", "Database user")
	RootCmd.PersistentFlags().StringVar(&dbPassword, "password", "shared_password", "Database password")
	RootCmd.PersistentFlags().StringVar(&dbName, "dbname", "shared_db", "Database name")
	RootCmd.PersistentFlags().StringVar(&dbSSLMode, "sslmode", "disable", "Database SSL mode")
}
