package xml

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestParsesExample(t *testing.T) {
	db, err := Parse("testdata/example-1.5.0.gramps")
	if err != nil {
		t.Fatalf("Failed to parse example: %s", err)
	}

	dir, err := ioutil.TempDir("/tmp", "xml-parse-test")
	if err != nil {
		t.Fatalf("Failed to create test dir: %s", err)
	}
	defer os.RemoveAll(dir)

	af := filepath.Join(dir, "actual-1.5.0.gramps")
	if err = db.Serialize(filepath.Join(dir, "actual-1.5.0.gramps")); err != nil {
		t.Fatalf("Failed to serialize db: %s", err)
	}

	ab, err := ioutil.ReadFile(af)
	if err != nil {
		t.Fatalf("Failed to read actual output: %s", err)
	}
	gf := "testdata/golden-1.5.0.gramps"
	gb, err := ioutil.ReadFile(gf)
	if err != nil {
		t.Fatalf("Failed to read golden output: %s", err)
	}
	if !bytes.Equal(ab, gb) {
		nf := "testdata/actual-1.5.0.gramps"
		os.Rename(af, nf)
		t.Errorf("Actual not equal to golden %s vs %s", nf, gf)
	}
}
