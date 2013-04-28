package xml

import (
	"bytes"
  "flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
  testDir = flag.String("dir", "testdata", "Directory to find test data")
)

func BenchmarkParse(b *testing.B) {
  b.StopTimer()
  fb, err := ioutil.ReadFile(filepath.Join(*testDir, "example-1.5.0.gramps"))
  if err != nil {
    b.Fatalf("Failed to read file: %s", err)
	}
  b.StartTimer()
  for i := 0; i < b.N; i++ {
    _, err = Parse(bytes.NewReader(fb))
    if err != nil {
      b.Fatalf("Failed to parse example: %s", err)
    }
  }
}

func TestParsesExample(t *testing.T) {
  f, err := os.Open(filepath.Join(*testDir, "example-1.5.0.gramps"))
  if err != nil {
    t.Fatalf("Failed to open file: %s", err)
	}
  db, err := Parse(f)
	if err != nil {
		t.Fatalf("Failed to parse example: %s", err)
	}

  tmpDir, err := ioutil.TempDir("/tmp", "xml-parse-test")
	if err != nil {
    t.Fatalf("Failed to create tmp dir: %s", err)
	}
  defer os.RemoveAll(tmpDir)

  af := filepath.Join(tmpDir, "actual-1.5.0.gramps")
  if err = db.Serialize(af); err != nil {
		t.Fatalf("Failed to serialize db: %s", err)
	}

	ab, err := ioutil.ReadFile(af)
	if err != nil {
		t.Fatalf("Failed to read actual output: %s", err)
	}
  gf := filepath.Join(*testDir, "golden-1.5.0.gramps")
	gb, err := ioutil.ReadFile(gf)
	if err != nil {
		t.Fatalf("Failed to read golden output: %s", err)
	}
	if !bytes.Equal(ab, gb) {
    nf := filepath.Join(*testDir, "/actual-1.5.0.gramps")
		os.Rename(af, nf)
		t.Errorf("Actual not equal to golden %s vs %s", nf, gf)
	}
}

func init() {
  flag.Parse()
}
