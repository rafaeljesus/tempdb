package tempdb

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"gopkg.in/redis.v5"
)

const (
	prefix = "tempdb:"
)

var (
	// ErrKeyRequired is returned when key is not passed as parameter in tempdb.Config.
	ErrKeyRequired = errors.New("creating a new tempdb requires a non-empty key")
	// ErrValueRequired is returned when value is not passed as parameter in tempdb.Config.
	ErrValueRequired = errors.New("creating a new tempdb requires a non-empty value")
)

type (
	// Tempdb stores an expiring (or non-expiring) key/value pair in Redis.
	Tempdb struct {
		*redis.Client
	}
	// Options carries the different variables to tune a newly started redis client,
	// it exposes the same configuration available from https://godoc.org/gopkg.in/redis.v5#Options go client.
	Options struct {
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
)

// New returns a new temp configured with the
// variables from the options parameter, or returning an non-nil err
// if an error occurred while creating redis client.
func New(o Options) (*Tempdb, error) {
	opts := newOptions(o)
	client := redis.NewClient(opts)
	return &Tempdb{client}, nil
}

// Insert a key/value pair with an optional expiration time.
func (t *Tempdb) Insert(key, value string, expires time.Duration) error {
	if len(key) == 0 {
		return ErrKeyRequired
	}

	if len(value) == 0 {
		return ErrValueRequired
	}

	k := fmt.Sprint(prefix, key)
	return t.Set(k, value, expires).Err()
}

// Find the value associated with the key.
func (t *Tempdb) Find(key string) (string, error) {
	if len(key) == 0 {
		return "", ErrKeyRequired
	}

	k := fmt.Sprint(prefix, key)
	value, err := t.Get(k).Result()
	if err != nil {
		return "", err
	}

	return value, t.Del(k).Err()
}

func newOptions(o Options) *redis.Options {
	opts := &redis.Options{}
	setAddr(opts, o.Addr)
	setPassword(opts, o.Password)
	setDB(opts, o.DB)
	setNetwork(opts, o.Network)
	setDialer(opts, o.Dialer)
	setMaxRetries(opts, o.MaxRetries)
	setDialTimeout(opts, o.DialTimeout)
	setReadTimeout(opts, o.ReadTimeout)
	setWriteTimeout(opts, o.WriteTimeout)
	setPoolSize(opts, o.PoolSize)
	setPoolTimeout(opts, o.PoolTimeout)
	setIdleTimeout(opts, o.IdleTimeout)
	setIdleCheckFrequency(opts, o.IdleCheckFrequency)
	setReadOnly(opts, o.ReadOnly)
	setTLSConfig(opts, o.TLSConfig)
	return opts
}

func setAddr(opts *redis.Options, addr string) {
	if len(addr) == 0 {
		addr = "localhost:6379"
	}
	opts.Addr = addr
}

func setPassword(opts *redis.Options, password string) {
	if len(password) != 0 {
		opts.Password = password
	}
}

func setDB(opts *redis.Options, db int) {
	if db != 0 {
		opts.DB = db
	}
}

func setNetwork(opts *redis.Options, network string) {
	if len(network) != 0 {
		opts.Network = network
	}
}

func setDialer(opts *redis.Options, dialer func() (net.Conn, error)) {
	if dialer != nil {
		opts.Dialer = dialer
	}
}

func setMaxRetries(opts *redis.Options, maxRetries int) {
	if maxRetries != 0 {
		opts.MaxRetries = maxRetries
	}
}

func setDialTimeout(opts *redis.Options, dialTimeout time.Duration) {
	if dialTimeout != 0 {
		opts.DialTimeout = dialTimeout
	}
}

func setReadTimeout(opts *redis.Options, readTimeout time.Duration) {
	if readTimeout != 0 {
		opts.ReadTimeout = readTimeout
	}
}

func setWriteTimeout(opts *redis.Options, writeTimeout time.Duration) {
	if writeTimeout != 0 {
		opts.WriteTimeout = writeTimeout
	}
}

func setPoolSize(opts *redis.Options, poolSize int) {
	if poolSize != 0 {
		opts.PoolSize = poolSize
	}
}

func setPoolTimeout(opts *redis.Options, poolTimeout time.Duration) {
	if poolTimeout != 0 {
		opts.PoolTimeout = poolTimeout
	}
}

func setIdleTimeout(opts *redis.Options, idleTimeout time.Duration) {
	if idleTimeout != 0 {
		opts.IdleTimeout = idleTimeout
	}
}

func setIdleCheckFrequency(opts *redis.Options, idleCheckFrequency time.Duration) {
	if idleCheckFrequency != 0 {
		opts.IdleCheckFrequency = idleCheckFrequency
	}
}

func setReadOnly(opts *redis.Options, readOnly bool) {
	if readOnly {
		opts.ReadOnly = readOnly
	}
}

func setTLSConfig(opts *redis.Options, tlsConfig *tls.Config) {
	if tlsConfig != nil {
		opts.TLSConfig = tlsConfig
	}
}
