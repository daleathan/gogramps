package xml

import (
  "testing"
)

func TestParsesExample(t *testing.T) {
  _, err := Parse("testdata/example-1.5.0.gramps")
  if err != nil {
    t.Errorf("Failed to parse example: %s", err)
  }
}
