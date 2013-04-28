package xml

import (
	"encoding/xml"
	"fmt"
	"testing"
)

const (
	xmlns = `xmlns="` + XMLNamespace + `"`
)

func TestGetDateString(t *testing.T) {
	cases := []struct{ xml, date string }{
		{fmt.Sprintf(`<dateval %s val="test"/>`, xmlns), "test"},
		{fmt.Sprintf(`<datestr %s val="test"/>`, xmlns), "test"},
		{fmt.Sprintf(`<daterange %s start="A" stop="B"/>`, xmlns),
			"between A and B"},
		{fmt.Sprintf(`<datespan %s start="C" stop="D"/>`, xmlns),
			"from C to D"},
	}
	for _, c := range cases {
		d := new(hasDate)
		if err := xml.Unmarshal([]byte("<test>"+c.xml+"</test>"), d); err != nil {
			t.Fatalf("Couldn't parse: [%s]: %s", c.xml, err)
		}
		if d.GetDateString() != c.date {
			t.Errorf("Expected [%s], got [%s", c.date, d.GetDateString())
		}
	}
}
