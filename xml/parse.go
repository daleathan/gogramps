package xml

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
  "sort"
)

func findUnparsedFields(v interface{}, path string, unparsed *map[string]bool) {
	value := reflect.ValueOf(v)
	fv := value.FieldByName("Unparsed")
	if !fv.IsValid() {
		panic(fmt.Sprintf("Unable to find Unparsed field in type: %s",
			value.Type()))
	}
	unparsedXML := fv.Interface().([]raw)

	fv = value.FieldByName("XMLName")
	if !fv.IsValid() {
		panic(fmt.Sprintf("Unable to find XMLName field in type: %s",
			value.Type()))
	}
	newPath := fv.Interface().(xml.Name).Local
	if path != "" {
		newPath = path + "." + newPath
	}

	for _, v := range unparsedXML {
		(*unparsed)[newPath+"."+v.XMLName.Local] = true
	}

	n := value.NumField()
	for i := 0; i < n; i++ {
		subField := value.Field(i)
		subType := value.Type().Field(i)
		if value.Type().Field(i).Name == "XMLName" {
			continue
		}
		if subField.Type().String() == "xml.raw" ||
			subField.Type().String() == "[]xml.raw" {
			continue
		}
		if subField.Kind() == reflect.Struct && !subType.Anonymous {
			findUnparsedFields(subField.Interface(), newPath, unparsed)
		}
		if subField.Kind() == reflect.Slice &&
			subField.Type().Elem().Kind() == reflect.Struct {
			n := subField.Len()
			for j := 0; j < n; j++ {
				findUnparsedFields(subField.Index(j).Interface(), newPath, unparsed)
			}
		}
	}
}

func fullyParsed(db Database) error {
	unparsedFields := make(map[string]bool)
	findUnparsedFields(db, "", &unparsedFields)
	if len(unparsedFields) > 0 {
		names := make([]string, 0, len(unparsedFields))
		for k := range unparsedFields {
      names = append(names, k)
		}
    sort.Strings(names)
		return fmt.Errorf("Unparsed fields: %s", names)
	}
	return nil
}

// Unmarshal a .gramps XML file in an arbitrary interface with XML tags.
// Useful for parsing a portion of a database.
func Unmarshal(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	unzipped, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer unzipped.Close()

	data, err := ioutil.ReadAll(unzipped)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

// Parse a .gramps XML file into a full Database. Returns an error if any
// portion of the XML was unparsed.
func Parse(filename string) (*Database, error) {
	var parsed Database
	if err := Unmarshal(filename, &parsed); err != nil {
		return nil, err
	}

	if err := fullyParsed(parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

// Serialize the Database back out into a .gramps XML file.
func (db Database) Serialize(filename string) error {
	data, err := xml.MarshalIndent(db, "", "\t")
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	zipped := gzip.NewWriter(f)
	if err != nil {
		return err
	}

	_, err = zipped.Write(data)
	zipped.Close()
	if err != nil {
		return err
	}

	return nil
}
