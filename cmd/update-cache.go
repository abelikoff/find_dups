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

var FullCacheUpdate bool

// updateCacheCmd represents the update-cache command

var updateCacheCmd = &cobra.Command{
	Use:   "update-cache",
	Short: "Update the cache file",
	Long: `Caches signatures for all files underneath a given directory. Also cleans entries from the cache
file that have been created more than 3 months ago`,
	Args: func(cmd *cobra.Command, args []string) error {
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
		UseCache = true
		VerbosityLevel, _ = cmd.Flags().GetCount("verbose")

		if VerbosityLevel > 0 {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		err := updateCache(args[0])

		if err != nil {
			logger.Fatal().Msgf("ERROR: %e", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCacheCmd)

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

	updateCacheCmd.PersistentFlags().BoolVarP(&FullCacheUpdate, "full", "f", false, "full update (default: incremental)")
	updateCacheCmd.Flags().StringVarP(&CacheFile, "cache_file", "", "",
		"cache file to use (default: ~/.find_dups.cache)")

}
