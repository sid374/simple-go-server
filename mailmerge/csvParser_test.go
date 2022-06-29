package mailmerge

import (
	"fmt"
	"os"
	"testing"

	"github.com/jedib0t/go-pretty/table"
)

func TestReadCsv(t *testing.T) {
	output, headers := ReadCSVFile("testing/m.csv")
	fmt.Printf("%v", output)

	pt := table.NewWriter()

	pt.SetOutputMirror(os.Stdout)
	headerRow := table.Row{}
	for _, col  := range headers {
		headerRow = append(headerRow, col)
	}

	pt.AppendHeader(headerRow)


	for _, row := range output {
		tableRow := table.Row{}
		for _, header := range headers {
			tableRow = append(tableRow, row[header])	
		}
		pt.AppendRow(tableRow)
	}
	pt.Render()
}