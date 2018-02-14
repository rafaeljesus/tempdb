package tempdb

import (
	"crypto/tls"
	"net"
	"testing"
	"time"
)

func TestTempdb(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T)
	}{
		{
			"new",
			testNew,
		},
		{
			"insert",
			testInsert,
		},
		{
			"find",
			testFind,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t)
		})
	}
}

func testNew(t *testing.T) {
	_, err := New(Options{
		Network:            "foo",
		Addr:               "localhost:6379",
		Dialer:             func() (c net.Conn, err error) { return },
		DB:                 1,
		Password:           "foo",
		MaxRetries:         1,
		DialTimeout:        time.Second * 5,
		ReadTimeout:        time.Second * 5,
		WriteTimeout:       time.Second * 5,
		PoolSize:           1,
		PoolTimeout:        time.Second * 5,
		IdleTimeout:        time.Second * 5,
		IdleCheckFrequency: time.Second * 5,
		ReadOnly:           true,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	if err != nil {
		t.Fatalf("expected to initialize tempdb: %s", err)
	}
}

func testInsert(t *testing.T) {
	temp, err := New(Options{})
	if err != nil {
		t.Fatalf("expected to initialize tempdb: %s", err)
	}

	cases := []struct {
		key, value string
		wantErr    bool
	}{
		{
			"key", "value", false,
		},
		{
			"", "value", true,
		},
		{
			"key", "", true,
		},
	}

	for _, c := range cases {
		err := temp.Insert(c.key, c.value, 0)
		if c.wantErr {
			if err == nil {
				t.Fatalf("unexpected insert return value: %v", err)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected insert return value: %v", err)
			}
		}
	}
}

func testFind(t *testing.T) {
	temp, err := New(Options{})
	if err != nil {
		t.Fatalf("expected to initialize tempdb: %v", err)
	}

	if err := temp.Insert("key", "value", 0); err != nil {
		t.Fatalf("expected to insert key/value: %v", err)
	}

	cases := []struct {
		key     string
		wantErr bool
	}{
		{
			"", true,
		},
		{
			"invalid_key", true,
		},
		{
			"key", false,
		},
	}

	for _, c := range cases {
		_, err := temp.Find(c.key)
		if c.wantErr {
			if err == nil {
				t.Fatalf("unexpected find return value")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected find return value")
			}
		}
	}
}
