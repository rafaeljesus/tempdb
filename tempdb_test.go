package tempdb

import (
	"testing"
)

func TestInsert(t *testing.T) {
	temp, err := NewTempdb(Options{})
	if err != nil {
		t.Errorf("Expected to initialize tempdb %s", err)
	}

	if err := temp.Insert("key", "value", 0); err != nil {
		t.Errorf("Expected to insert key/value %s", err)
	}
}

func TestGet(t *testing.T) {
	temp, err := NewTempdb(Options{})
	if err != nil {
		t.Errorf("Expected to initialize tempdb %s", err)
	}

	if err := temp.Insert("key", "value", 0); err != nil {
		t.Errorf("Expected to insert key/value %s", err)
	}

	value, err := temp.Find("key")
	if err != nil {
		t.Errorf("Expected get value %s", err)
	}

	if value != "value" {
		t.Errorf("Expected value to be eq %s", value)
	}
}
