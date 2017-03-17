package tempdb

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gopkg.in/redis.v5"
	"net"
	"time"
)

var (
	// ErrKeyRequired is returned when key is not passed as parameter in tempdb.Config.
	ErrKeyRequired = errors.New("creating a new tempdb requires a non-empty key")
	// ErrValueRequired is returned when value is not passed as parameter in tempdb.Config.
	ErrValueRequired = errors.New("creating a new tempdb requires a non-empty value")
)

// Tempdb stores an expiring (or non-expiring) key/value pair in Redis.
type Tempdb interface {
	Insert(key, value string, expires time.Duration) (err error)
	Find(key string) (value string, err error)
}

// Options carries the different variables to tune a newly started redis client,
// it exposes the same configuration available from https://godoc.org/gopkg.in/redis.v5#Options go client.
type Options struct {
	Network            string
	Addr               string
	Dialer             func() (net.Conn, error)
	Password           string
	DB                 int
	MaxRetries         int
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	ReadOnly           bool
	TLSConfig          *tls.Config
}

type temp struct {
	*redis.Client
}

// NewTempdb returns a new temp configured with the
// variables from the options parameter, or returning an non-nil err
// if an error ocurred while creating redis client.
func NewTempdb(o Options) (tempdb Tempdb, err error) {
	options := newOptions(o)
	client := redis.NewClient(options)
	tempdb = &temp{client}

	return
}

// Insert a key/value pair with an optional expiration time.
func (t *temp) Insert(key, value string, expires time.Duration) (err error) {
	if len(key) == 0 {
		err = ErrKeyRequired
		return
	}

	if len(value) == 0 {
		err = ErrValueRequired
		return
	}

	k := fmt.Sprint("tempDB:", key)
	err = t.Set(k, value, expires).Err()

	return
}

// Find the value associated with the key.
func (t *temp) Find(key string) (value string, err error) {
	if len(key) == 0 {
		err = ErrKeyRequired
		return
	}

	k := fmt.Sprint("tempDB:", key)
	value, err = t.Get(k).Result()
	if err != nil {
		return
	}

	err = t.Del(k).Err()

	return
}

func newOptions(o Options) (options *redis.Options) {
	options = &redis.Options{}

	if len(o.Addr) != 0 {
		options.Addr = o.Addr
	}

	if len(o.Password) != 0 {
		options.Password = o.Password
	}

	if o.DB != 0 {
		options.DB = o.DB
	}

	if len(o.Network) != 0 {
		options.Network = o.Network
	}

	if o.Dialer != nil {
		options.Dialer = o.Dialer
	}

	if o.MaxRetries != 0 {
		options.MaxRetries = o.MaxRetries
	}

	if o.DialTimeout != 0 {
		options.DialTimeout = o.DialTimeout
	}

	if o.ReadTimeout != 0 {
		options.ReadTimeout = o.ReadTimeout
	}

	if o.WriteTimeout != 0 {
		options.WriteTimeout = o.WriteTimeout
	}

	if o.PoolSize != 0 {
		options.PoolSize = o.PoolSize
	}

	if o.PoolTimeout != 0 {
		options.PoolTimeout = o.PoolTimeout
	}

	if o.IdleTimeout != 0 {
		options.IdleTimeout = o.IdleTimeout
	}

	if o.IdleCheckFrequency != 0 {
		options.IdleCheckFrequency = o.IdleCheckFrequency
	}

	if o.ReadOnly {
		options.ReadOnly = o.ReadOnly
	}

	if o.TLSConfig != nil {
		options.TLSConfig = o.TLSConfig
	}

	return
}
