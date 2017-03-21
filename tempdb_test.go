package tempdb

import (
	"testing"
)

func TestInsert(t *testing.T) {
	temp, err := New(Options{})
	if err != nil {
		t.Errorf("Expected to initialize tempdb %s", err)
	}

	if err := temp.Insert("key", "value", 0); err != nil {
		t.Errorf("Expected to insert key/value %s", err)
	}

	if err := temp.Insert("", "value", 0); err == nil {
		t.Fail()
	}

	if err := temp.Insert("key", "", 0); err == nil {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	temp, err := New(Options{})
	if err != nil {
		t.Errorf("Expected to initialize tempdb %s", err)
	}

	if err := temp.Insert("key", "value", 0); err != nil {
		t.Errorf("Expected to insert key/value %s", err)
	}

	_, err = temp.Find("")
	if err == nil {
		t.Fail()
	}

	_, err = temp.Find("invalid_key")
	if err == nil {
		t.Fail()
	}

	value, err := temp.Find("key")
	if err != nil {
		t.Errorf("Expected get value %s", err)
	}

	if value != "value" {
		t.Errorf("Expected value to be eq %s", value)
	}
}
