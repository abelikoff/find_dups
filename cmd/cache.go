// Copyright © 2024 Alexander L. Belikoff (alexander@belikoff.net)

// Caching facility for find_dups.

package cmd

import (
	"fmt"
	"log"

	"github.com/colinmarc/cdb"
)

func getCachedFileInfo(path string) (FileInfo, bool) {
	return nil, false
}

func cacheFileInfo(file_info FileInfo) bool {
	return false
}

func main() {
	writer, err := cdb.Create("/tmp/example.cdb")
	if err != nil {
		log.Fatal(err)
	}

	// Write some key/value pairs to the database.
	writer.Put([]byte("Alice"), []byte("Practice"))
	writer.Put([]byte("Bob"), []byte("Hope"))
	writer.Put([]byte("Charlie"), []byte("Horse"))

	// Freeze the database, and open it for reads.
	db, err := writer.Freeze()
	if err != nil {
		log.Fatal(err)
	}

	// Fetch a value.
	v, err := db.Get([]byte("Alice"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(v))
}
