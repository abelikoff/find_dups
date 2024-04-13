// Copyright © 2024 Alexander L. Belikoff (alexander@belikoff.net)

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "find_dups <directory>",
	Short: "List all identical files underneath the directory",
	Long: `Files are considered	identical if they have the same size
and SHA-1 signature.`,

	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		path := args[0]
		fileInfo, err := os.Stat(path)

		if err != nil || !fileInfo.IsDir() {
			return fmt.Errorf("not a directory: %s", path)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := groupBySize(args[0])

		if err != nil {
			log.Fatal(err)
		}

		groupBySignature()
		showDuplicates()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.find_dups.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
