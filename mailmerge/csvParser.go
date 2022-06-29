package mailmerge

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func ReadCSVFile(filePath string) ([]map[string]string, []string) {
	f, err := os.Open(filePath)
    if err != nil {
        log.Fatal(err)
    }

    // remember to close the file at the end of the program
    defer f.Close()

    // read csv values using csv.Reader
    csvReader := csv.NewReader(f)
    data, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

	if len(data) == 0 {
		log.Fatal("Empty CSV found")
	}
	headers := []string {}

	// Parse Headers
	for _, col := range data[0] {
		headers = append(headers, strings.Trim(col, " "))
	}

	output_csv := []map[string]string{}

	for _, row := range data[1:] {
		m := map[string]string{}
		output_csv = append(output_csv, m)
		for i, col := range row {
			m[headers[i]] = strings.Trim(col, " ")
		}
	}

	return output_csv, headers
}