## gFind

A simple utility for finding files in a network / slow directory using CSV indexes for performance

Note: 
- Searches are case-insensitive
- All search tokens must match for a file to be reported

Build:

```shell
go build -o gFind
```

Usage:

```shell
./gFind -dir=[string] -refresh=[bool]* -minSize=[int64]* [search_tokens]

* - optional
```

Example:

```shell
./gFind -dir=/Volumes/Storage/ebooks italian
```