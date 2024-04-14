// Copyright © 2024 Alexander L. Belikoff (alexander@belikoff.net)

// http://github.com/abelikoff/find_dups

// Caching facility for find_dups.
//
// For each file (of those already similar by size) we store signature as well as the file size.
// When fetching the cached value, we make sure the size hasn't changed.
//
// We also store the cache write timespamp for cache cleanups (currently not implemented).

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/graygnuorg/go-gdbm"
)

// const ExpirationHours float64 = 24 * 30

var num_cache_calls, num_cache_hits, num_cache_mismatches, num_cache_writes int // cache statistics

type CachedEntry struct {
	Signature string
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

func getCachedFileInfo(info *FileInfo) (*CachedEntry, error) {
	num_cache_calls++

	if db == nil {
		if err := openCache(); err != nil {
			return nil, fmt.Errorf("ERROR: failed to open cache file: %e", err)
		}
	}

	value, err := db.Fetch([]byte(info.Path))

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

	if len(tokens) != 3 {
		return nil, fmt.Errorf("ERROR: malformed cache entry: '%s' -> '%s'", info.Path, serialized)
	}

	entry.Signature = tokens[0]

	// make sure cached entry has not expired

	/*save_time, err := time.Parse(time.RFC3339, tokens[2])

	if err != nil {
		return nil, fmt.Errorf("ERROR: malformed cache time for '%s': '%s': %e", path, tokens[2], err)
	}

	now := time.Now()
	time_diff := now.Sub(save_time)

	if time_diff.Hours() >= ExpirationHours {
		num_cache_expires++
		db.Delete([]byte(path))
		return nil, nil
	} */

	// get file size

	cached_size, err := strconv.ParseInt(tokens[1], 10, 64)

	if err != nil {
		return nil, fmt.Errorf("ERROR: malformed cached size for '%s': '%s': %e", info.Path, tokens[1], err)
	}

	if cached_size != info.Size {
		logger.Debug().Msgf("cache mismatch for '%s': cached size: %d, real size: %d",
			info.Path, cached_size, info.Size)
		num_cache_mismatches++
		db.Delete([]byte(info.Path))
		return nil, nil
	}

	if VerbosityLevel > 1 {
		logger.Debug().Msgf("cache hit: '%s' -> '%s'", info.Path, entry.Signature)
	}

	num_cache_hits++
	return &entry, nil
}

func cacheFileInfo(info *FileInfo, data *CachedEntry) error {
	serialized := fmt.Sprintf("%s|%d|%s", data.Signature, info.Size, time.Now().Format(time.RFC3339))
	//serialized := fmt.Sprintf("%s|%s", data.Signature, time.Now().Format(time.RFC3339))

	if VerbosityLevel > 1 {
		logger.Debug().Msgf("cache write: '%s' -> '%s'", info.Path, serialized)
	}

	err := db.Store([]byte(info.Path), []byte(serialized), true)

	if err != nil {
		return fmt.Errorf("ERROR: failed to cache '%s': %e", info.Path, err)
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

	logger.Debug().Msgf("cache stats: %d total, %d hits (%.1f%% hit ratio), %d misses (%d mismatches), %d writes",
		num_cache_calls, num_cache_hits, hit_ratio, num_cache_calls-num_cache_hits,
		num_cache_mismatches, num_cache_writes)
}
