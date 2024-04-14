// Copyright © 2024 Alexander L. Belikoff (alexander@belikoff.net)

// http://github.com/abelikoff/find_dups

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var logger zerolog.Logger
var VerbosityLevel int
var UseCache bool
var CacheFile string

// base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "find_dups  <directory> ",
	Short: "List all identical files underneath the directory",
	Long: `find_dups lists all identical files underneath the directory.

Files are considered identical if they have the same size and SHA-1 signature.
File signatures can be cached for faster operation.

`,

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
		VerbosityLevel, _ = cmd.Flags().GetCount("verbose")

		if VerbosityLevel > 0 {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		logger.Debug().Msg("grouping files by size...")
		err := groupBySize(args[0])

		if err != nil {
			logger.Fatal().Msgf("ERROR: %e", err)
		}

		logger.Debug().Msg("grouping files by signature...")
		groupBySignature()
		logger.Debug().Msg("generating report...")
		showDuplicates()
		outputCacheStats()
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

	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	rootCmd.PersistentFlags().BoolVarP(&UseCache, "cache", "C", false, "enable caching")
	rootCmd.Flags().StringVarP(&CacheFile, "cache_file", "", "",
		"cache file to use (default: ~/.find_dups.cache)")
	rootCmd.PersistentFlags().CountP("verbose", "v", "verbosity level")
}
