package mailmerge

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
)

type Template struct {
	varName string
	replacedValue string
	startPosition int
	endPosition int
}

var licenseKey = ``

func getStructFromSubmatchIndex(submatch []int, s string, mapping map[string]string) Template {
	varName := s[submatch[2]:submatch[3]]
	return Template{
		varName: varName,
		replacedValue: mapping[varName],
		startPosition: submatch[0],
		endPosition: submatch[1],
	}
}

func GetTemplates(s string, mapping map[string]string) (templates []Template) {
	re := regexp.MustCompile(`<([^>]*)>`)
	matches := re.FindAllStringSubmatchIndex(s, -1)

	templates = []Template{}
	for _, subarray := range matches {
		temp := getStructFromSubmatchIndex(subarray, s, mapping)
		templates = append(templates, temp)
	}
	return
}

func Merge(s string, mapping map[string]string) string {
	templates := GetTemplates(s, mapping)
	sort.Slice(templates, func(i, j int) bool { return templates[i].startPosition < templates[j].startPosition})
	
	template_idx := 0
	string_idx := 0
	final := ""
	for string_idx < len(s) {
		if template_idx < len(templates) && string_idx == templates[template_idx].startPosition {
			template := templates[template_idx]
			orig_index := 0
			for orig_index < len(template.replacedValue) {
				final = final + string(template.replacedValue[orig_index])
				orig_index += 1
			}
			string_idx = template.endPosition
			template_idx += 1
		} else {
			final = final + string(s[string_idx])
			string_idx += 1
		}
	}
	
	return final 
}

func CreateMergeMap(
	firstName string, 
	lastName string, 
	title string, 
	firmNameFull string, 
	firmName string, 
	address1 string,
	address2 string) map[string]string {
		m := map[string]string {
			"FirstName": firstName,
			"LastName": lastName,
			"Title": title,
			"FirmNameFull": firmNameFull,
			"FirmName": firmName,
			"Address1": address1,
			"Address2": address2,
		}
		for k, val := range m {
			if len(val) < 2 {
				fmt.Printf("Empty value found for key %s, firm: %s\n", k, m["FirmName"])
			}
		}
		return m
	}

func WriteDoc(s string, savePath string) {
	status, err := license.GetMeteredState()
	if err != nil {
		err := license.SetMeteredKey(licenseKey)
		if err != nil {
			fmt.Printf("Error loading license: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("License Status: %v\n", status)
	}
    

	doc := document.New()
	defer  doc.Close()

	para := doc.AddParagraph()
	run := para.AddRun()
	run.Properties().SetFontFamily("Times New Roman")
	run.Properties().SetSize(12)
	run.AddText(s)
	section := doc.BodySection()
	section.SetPageSizeAndOrientation(measurement.Inch*8.5, measurement.Inch*11, wml.ST_PageOrientationPortrait)

	err = doc.SaveToFile(savePath)
	if err != nil {
		fmt.Printf("Something went wrong writing doc %s\n", savePath)
	} else {
		fmt.Printf("Done writing doc! %s\n", savePath)
	}
}

func RowToMergeMap(row map[string]string) map[string]string {
	required := mapset.NewSet("Firm", "Title", "Name", "Address1", "Address2")
	for r := range required.Iterator().C {
		_, ok := row[r]
		if !ok {
			log.Fatal(fmt.Sprintf("Header %s not found", r))
		}
	}

	names := strings.Split(row["Name"], " ")
	lastName := names[len(names)-1]
	var title string
	var firstName string
	switch {
	case names[0] == "Mr." || names[0] == "Mr":
		title = "Mr."
	case names[0] == "Ms." || names[0] == "Ms":
		title = "Ms."
	}

	if title == "" {
		firstName = strings.Join(names[:len(names)-1], " ")
		title = row["Title"]
	} else {
		firstName = strings.Join(names[1:len(names)-1], " ")
	}

	if title == "" {
		fmt.Printf("WARNING: Title missing for firm %s", row["Firm"])
	}

	var fullFirmName string

	firmName := row["Firm"]
	lastWordIndex := strings.LastIndex(firmName, " ")
	if lastWordIndex != -1 && firmName[lastWordIndex:] == "LLP" {
		fullFirmName = firmName
		firmName = firmName[:lastWordIndex]
	} else {
		fullFirmName = firmName + " LLP"
	}

	return CreateMergeMap(
		firstName,
		lastName,
		title,
		fullFirmName,
		firmName,
		row["Address1"],
		row["Address2"],
	)
}



func CreateMarieMerge() {
	input_csv, _:= ReadCSVFile("testing/m.csv")

	for i, row := range input_csv {
		if i == 1 {
			return
		}
		mergeMap := RowToMergeMap(row)
		mergedEmail := Merge(EmailContents, mergeMap)
		WriteDoc(mergedEmail, fmt.Sprintf("output/%s.docx", mergeMap["FirmName"]))
	}
}


