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

### Use caching

`find_dups` can cache signatures for files it processes to facilitate faster
results upon repeated use.

```shell
    find_dups -C ~/mydir
```

### Cache clean-up

Remove all cache entries older than 3 months:

```shell
    find_dups cleanCache
```
