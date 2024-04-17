# find_dups

`find_dups` is a tool to find duplicate files underneath a given directory.

## Installation

```shell
    go install github.com/abelikoff/find_dups@latest
```

## Usage

### Basic usage

Find all identical files in a given directory:

```shell
    find_dups ~/mydir
```

### Caching

`find_dups` can cache signatures for files it processes to facilitate faster
results upon repeated use. Cache can be built in advance (which could be slow as it will
calculate signatures for all files in the directory) on on the fly (which will only
cache signatures for the files for which they are computed)

To rebuild the cache (_slow!_):

```shell
    find_dups update-cache
```

By default, cache rebuild is incremental (files that already have cached signatures are
not touched). Passing `-f` flag makes `find_dups` re-cache all files in the specified
directory.

To use cached signatures:

```shell
    find_dups -C ~/mydir
```
