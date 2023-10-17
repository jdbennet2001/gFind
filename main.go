package main

import (
	"errors"
	"flag"
	"fmt"
	"gFind/data"
	"gFind/defs"
	"os"
	"path/filepath"
	"strings"
)

/*
 Version of find that uses a cached CSV file to map network drives for fast scanning
 Usage: gFine -dir=string -refresh=bool*  [string]
 */
func main() {

	// Define command line flags
	tgtDirectoryPtr := flag.String("dir", "default", "a string")
	refreshPtr := flag.Bool("refresh", false, "a bool")
	minSizePtr := flag.Int64("minSize", 32*1024, "an int64")

	// Parse the command line flags
	flag.Parse()

	// Get the command line tokens after feature flags
	searchTokens := flag.Args()

	/* -- Sanity Checks -- */

	// Input parameters
	if *tgtDirectoryPtr == "default" || len(searchTokens) == 0{
		fmt.Println("Usage: gFine -dir=[string] -refresh=[bool]* [searchTokens]")
		return
	}


	/* --------------- */

	// Need to refresh.. go for it..
	if *refreshPtr == true || !data.HasCSV(*tgtDirectoryPtr) {
		fmt.Println("Scanning input directory")
		files, err := Walk(*tgtDirectoryPtr, *minSizePtr)
		if err != nil {
			fmt.Println("Input directory not reachable: ", *tgtDirectoryPtr)
			return
		}
		data.WriteCSV(files, *tgtDirectoryPtr)
	}

	// Load data
	csv := data.ReadCSV(*tgtDirectoryPtr)

	for _, row := range(csv){
		if match(searchTokens, row.Location){
			rel, _ := filepath.Rel(*tgtDirectoryPtr, row.Location)
			fmt.Println(rel, " : ", formatBytes(row.Size))
		}
	}

	return

}


/*
 Return a map of <string: name> : <size: int64>
 for each comic in the root directory
*/
func Walk(rootpath string, minSize int64)([]defs.FileInfo, error) {

	// Input directory exists
	if _, err := os.Stat(rootpath); os.IsNotExist(err) {
		return nil, errors.New("Input directory does not exist")
	}

	var results []defs.FileInfo

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {

		println("Scanning ", path)

		size := info.Size()

		c := defs.FileInfo{
			Location: path,
			Size: size,
		}

		// Drop thumbnails
		if size > minSize{
			results = append(results, c)
		}

		return nil	// No error...
	})

	if err != nil {
		fmt.Printf("Walk error [%v]\n", err)
	}

	return results, nil

}

func match(searchTokens []string, str string) bool{

	for _, token := range(searchTokens){
		baseToken := filepath.Base(str)
		baseLower := strings.ToLower(baseToken)
		tokenLower := strings.ToLower(token)
		if ! strings.Contains(baseLower, tokenLower){
			return false
		}
	}

	return true
}

func formatBytes(count int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
		TB = 1 << 40
	)

	switch {
	case count < KB:
		return fmt.Sprintf("%d B", count)
	case count < MB:
		return fmt.Sprintf("%.2f KB", float64(count)/float64(KB))
	case count < GB:
		return fmt.Sprintf("%.2f MB", float64(count)/float64(MB))
	case count < TB:
		return fmt.Sprintf("%.2f GB", float64(count)/float64(GB))
	default:
		return fmt.Sprintf("%.2f TB", float64(count)/float64(TB))
	}
}