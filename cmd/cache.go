// Copyright © 2024 Alexander L. Belikoff (alexander@belikoff.net)

// http://github.com/abelikoff/find_dups

// Caching facility for find_dups.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/graygnuorg/go-gdbm"
)

const ExpirationHours float64 = 24

var num_cache_calls, num_cache_hits, num_cache_expires, num_cache_writes int // cache statistics

type CachedEntry struct {
	Signature string
	//Size      int64
}

var db *gdbm.Database

func openCache() error {
	if CacheFile == "" {
		home, _ := os.UserHomeDir()
		CacheFile = home + string(os.PathSeparator) + ".find_dups.cache"
	}

	var err error
	db, err = gdbm.Open(CacheFile, gdbm.ModeWrcreat)

	if err != nil {
		return fmt.Errorf("ERROR: failed to open cache file: %e", err)
	}

	return nil
}

func getCachedFileInfo(path string) (*CachedEntry, error) {
	num_cache_calls++

	if db == nil {
		if err := openCache(); err != nil {
			return nil, fmt.Errorf("ERROR: failed to open cache file: %e", err)
		}
	}

	value, err := db.Fetch([]byte(path))

	if err != nil {
		if errors.Is(err, gdbm.ErrItemNotFound) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("ERROR: GDBM error: %e", err)
		}
	}

	var serialized string = string(value)
	var entry CachedEntry

	tokens := strings.Split(serialized, "|")

	if len(tokens) != 2 {
		return nil, fmt.Errorf("ERROR: malformed cache entry: '%s' -> '%s'", path, serialized)
	}

	entry.Signature = tokens[0]

	// make sure cached entry has not expired

	save_time, err := time.Parse(time.RFC3339, tokens[1])

	if err != nil {
		return nil, fmt.Errorf("ERROR: malformed cache time for '%s': '%s': %e", path, tokens[1], err)
	}

	now := time.Now()
	time_diff := now.Sub(save_time)

	if time_diff.Hours() >= ExpirationHours {
		num_cache_expires++
		db.Delete([]byte(path))
		return nil, nil
	}

	// get file size

	/*entry.Size, err = strconv.ParseInt(tokens[2], 10, 64)

	if err != nil {
		return nil, fmt.Errorf("ERROR: malformed cached size for '%s': '%s': %e", path, tokens[2], err)
	}*/

	if VerbosityLevel > 1 {
		logger.Debug().Msgf("cache hit: '%s' -> '%s'", path, entry.Signature)
	}

	num_cache_hits++
	return &entry, nil
}

func cacheFileInfo(path string, data *CachedEntry) error {
	//serialized := fmt.Sprintf("%s|%d|%s", data.Signature, data.Size, time.Now().Format(time.RFC3339))
	serialized := fmt.Sprintf("%s|%s", data.Signature, time.Now().Format(time.RFC3339))

	if VerbosityLevel > 1 {
		logger.Debug().Msgf("cache write: '%s' -> '%s'", path, serialized)
	}

	err := db.Store([]byte(path), []byte(serialized), true)

	if err != nil {
		return fmt.Errorf("ERROR: failed to cache '%s': %e", path, err)
	}

	num_cache_writes++
	return nil
}

func outputCacheStats() {
	if VerbosityLevel < 1 {
		return
	}

	var hit_ratio float32 = 0.0

	if num_cache_hits > 0 {
		hit_ratio = float32(num_cache_hits) / float32(num_cache_calls) * 100
	}

	logger.Debug().Msgf("cache stats: %d total, %d hits (%.1f%% hit ratio), %d misses (%d expired), %d writes",
		num_cache_calls, num_cache_hits, hit_ratio, num_cache_calls-num_cache_hits,
		num_cache_expires, num_cache_writes)
}
