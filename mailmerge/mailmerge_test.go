package mailmerge

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestMerge(t *testing.T) {
//	s := "Dear <Title> <FirstName> <LastName>,\n How are you?"
	s := EmailContents
	m := map[string]string {
		"LastName": "Bid",
		"FirstName": "Sid",
		"Title": "Mr.",
		"FirmName": "Sid & Bid",
		"FirmNameFull": "Sid & Bid LLP",
		"Address1": "649 Arkansas Street",
		"Address2": "San Francisco, CA, 94107",
	}
	merged := Merge(s, m)
	fmt.Println(merged)
}

func TestDocCreation(t *testing.T) {
	WriteDoc(EmailContents, "testing/out.docx")
}

func TestMarieMerge(t *testing.T) {
	CreateMarieMerge()
}

func PrettyPrint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	var ret string
	if err == nil {
		ret = fmt.Sprintf(string(b))
	}
	return ret
}

func TestRowToMergeMap(t *testing.T) {
	out, _ := ReadCSVFile("testing/m_reordered_cols.csv")

	for i := 0; i < 5; i++ {
		mergeMap := RowToMergeMap(out[i])
		log.Printf(PrettyPrint(mergeMap))
	}
	
}