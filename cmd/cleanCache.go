// Copyright © 2024 Alexander L. Belikoff (alexander@belikoff.net)

// http://github.com/abelikoff/find_dups

package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// cleanCacheCmd represents the cleanCache command
var cleanCacheCmd = &cobra.Command{
	Use:   "clean-cache",
	Short: "Clean the cache file",
	Long: `Clean entries from the cache file that have been created more than
	3 months ago`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(0)(cmd, args); err != nil {
			return err
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

		err := cleanCache()

		if err != nil {
			logger.Fatal().Msgf("ERROR: %e", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCacheCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCacheCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCacheCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	cleanCacheCmd.Flags().StringVarP(&CacheFile, "cache_file", "", "",
		"cache file to use (default: ~/.find_dups.cache)")

}
