package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/templecloud/memoria-server/internal/memoria/boot"
)

var rootCmd = &cobra.Command{
	Use:   "memoria",
	Short: "Memoria is a memorization helper application.",
	Long:  `Memoria is a memorization helper application.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'memoria help' to find out how to use Memoria.")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Memoria",
	Long:  `Print the version number of Memoria`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// Display the Memoria version.
		fmt.Println("v0.0.1")
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Memoria server.",
	Long:  `Start the Memoria server.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadInitConfig()
		// Start the Memoria API server.
		boot.Start(cfg)
	},
}

func configure() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
}

// Execute run the Memoria root command.
func run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	configure()
	run()
}
