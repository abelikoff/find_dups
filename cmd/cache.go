// Copyright © 2024 Alexander L. Belikoff (alexander@belikoff.net)

// Caching facility for find_dups.

package cmd

import (
	"fmt"
	"strings"
	"time"
)

const ExpirationHours float64 = 24

type CachedEntry struct {
	Signature string
	//Size      int64
}

// var db *cdb.CDB

/* func openCache() error {
	home, _ := os.UserHomeDir()
	cache_file := home + string(os.PathSeparator) + ".find_dups.cache"
	var err error
	db, err = cdb.Open(cache_file)
	return err
	return nil
} */

func getCachedFileInfo(path string) (*CachedEntry, error) {
	// value, err := db.Get([]byte(path))

	var serialized string = ""
	var entry CachedEntry

	tokens := strings.Split(serialized, "|")

	if len(tokens) != 3 {
		return nil, fmt.Errorf("ERROR: malformed cache entry: '%s' -> '%s'", path, serialized)
	}

	entry.Signature = tokens[0]

	// make sure cached entry has not expired

	save_time, err := time.Parse(time.RFC3339, tokens[2])

	if err != nil {
		return nil, fmt.Errorf("ERROR: malformed cache time for '%s': '%s': %e", path, tokens[2], err)
	}

	now := time.Now()
	time_diff := now.Sub(save_time)

	if time_diff.Hours() >= ExpirationHours {
		// TODO: delete entry
		return nil, nil
	}

	// get file size

	/*entry.Size, err = strconv.ParseInt(tokens[1], 10, 64)

	if err != nil {
		return nil, fmt.Errorf("ERROR: malformed cached size for '%s': '%s': %e", path, tokens[1], err)
	}*/

	return &entry, nil
}

func cacheFileInfo(path string, data *CachedEntry) error {
	//serialized := fmt.Sprintf("%s|%d|%s", data.Signature, data.Size, time.Now().Format(time.RFC3339))
	serialized := fmt.Sprintf("%s|%s", data.Signature, time.Now().Format(time.RFC3339))
	fmt.Printf("*** cache put: '%s' -> '%s'\n", path, serialized)
	return nil
}
