package data

import (
"cacheGrep/defs"
"encoding/csv"
	"fmt"
	"log"
"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
 Write CSV files to disk
 @ return number of rows written, error
*/

func WriteCSV(files []defs.FileInfo, tgtPath string) error{

	csvName := csvName(tgtPath)

	csvFile, err := os.Create(csvName)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return err
	}

	defer  csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	// Write header
	row := []string{"Location", "Size"}
	csvwriter.Write(row)

	for _, c := range files{
		row := []string{c.Location, strconv.Itoa(int(c.Size))}
		csvwriter.Write(row)
	}

	defer csvwriter.Flush()

	return nil

}

func ReadCSV(tgtPath string) []defs.FileInfo{

	csv_location := csvName(tgtPath)

	// Pointer to a map[string]int64
	results := []defs.FileInfo{}

	f, err := os.Open(csv_location)
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(f)
	defer f.Close()

	//Read the header
	header, _ := csvReader.Read()
	fmt.Println("Header rows are: ", header)

	// Read data
	rows, _ := csvReader.ReadAll()
	for _,s := range rows{
		location := s[0]
		i64, _ := strconv.ParseInt(s[1], 10, 0)
		row := defs.FileInfo{location, i64}
		results = append(results, row)
	}

	return results
}

// Check if CSV data exists
func HasCSV(tgtPath string) bool{
	csvName := csvName(tgtPath)

	if _, err := os.Stat(csvName); os.IsNotExist(err) {
		return false
	}

	return true
}

// Translate a full qualified path into a single string (CSV) file
func csvName(fullPath string) string{

	// Replace directory separators with underscores
	filenameWithoutSeparators := strings.ReplaceAll(fullPath, string(filepath.Separator), "_")

	csvFileName := filenameWithoutSeparators +".csv"

	return csvFileName
}

