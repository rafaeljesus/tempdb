package tempdb

import (
	"crypto/tls"
	"errors"
	"fmt"
	d "github.com/tj/go-debug"
	"gopkg.in/redis.v5"
	"net"
	"time"
)

var (
	// ErrKeyRequired is returned when key is not passed as parameter in tempdb.Config.
	ErrKeyRequired = errors.New("creating a new tempdb requires a non-empty key")
	// ErrValueRequired is returned when value is not passed as parameter in tempdb.Config.
	ErrValueRequired = errors.New("creating a new tempdb requires a non-empty value")

	debug = d.Debug("single")
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

// New returns a new temp configured with the
// variables from the options parameter, or returning an non-nil err
// if an error occurred while creating redis client.
func New(o Options) (tempdb Tempdb, err error) {
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

	debug("Saved %s/%s. Expiring in %d seconds", k, value, expires)

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

	debug("Deleted %s/%s.", k, value)

	return
}

func newOptions(o Options) (options *redis.Options) {
	options = &redis.Options{}

	setAddr(options, o.Addr)
	setPassword(options, o.Password)
	setDB(options, o.DB)
	setNetwork(options, o.Network)
	setDialer(options, o.Dialer)
	setMaxRetries(options, o.MaxRetries)
	setDialTimeout(options, o.DialTimeout)
	setReadTimeout(options, o.ReadTimeout)
	setWriteTimeout(options, o.WriteTimeout)
	setPoolSize(options, o.PoolSize)
	setPoolTimeout(options, o.PoolTimeout)
	setIdleTimeout(options, o.IdleTimeout)
	setIdleCheckFrequency(options, o.IdleCheckFrequency)
	setReadOnly(options, o.ReadOnly)
	setTLSConfig(options, o.TLSConfig)

	return
}

func setAddr(options *redis.Options, addr string) {
	if len(addr) == 0 {
		addr = "localhost:6379"
	}
	options.Addr = addr
}

func setPassword(options *redis.Options, password string) {
	if len(password) != 0 {
		options.Password = password
	}
}

func setDB(options *redis.Options, db int) {
	if db != 0 {
		options.DB = db
	}
}

func setNetwork(options *redis.Options, network string) {
	if len(network) != 0 {
		options.Network = network
	}
}

func setDialer(options *redis.Options, dialer func() (net.Conn, error)) {
	if dialer != nil {
		options.Dialer = dialer
	}
}

func setMaxRetries(options *redis.Options, maxRetries int) {
	if maxRetries != 0 {
		options.MaxRetries = maxRetries
	}
}

func setDialTimeout(options *redis.Options, dialTimeout time.Duration) {
	if dialTimeout != 0 {
		options.DialTimeout = dialTimeout
	}
}

func setReadTimeout(options *redis.Options, readTimeout time.Duration) {
	if readTimeout != 0 {
		options.ReadTimeout = readTimeout
	}
}

func setWriteTimeout(options *redis.Options, writeTimeout time.Duration) {
	if writeTimeout != 0 {
		options.WriteTimeout = writeTimeout
	}
}

func setPoolSize(options *redis.Options, poolSize int) {
	if poolSize != 0 {
		options.PoolSize = poolSize
	}
}

func setPoolTimeout(options *redis.Options, poolTimeout time.Duration) {
	if poolTimeout != 0 {
		options.PoolTimeout = poolTimeout
	}
}

func setIdleTimeout(options *redis.Options, idleTimeout time.Duration) {
	if idleTimeout != 0 {
		options.IdleTimeout = idleTimeout
	}
}

func setIdleCheckFrequency(options *redis.Options, idleCheckFrequency time.Duration) {
	if idleCheckFrequency != 0 {
		options.IdleCheckFrequency = idleCheckFrequency
	}
}

func setReadOnly(options *redis.Options, readOnly bool) {
	if readOnly {
		options.ReadOnly = readOnly
	}
}

func setTLSConfig(options *redis.Options, tlsConfig *tls.Config) {
	if tlsConfig != nil {
		options.TLSConfig = tlsConfig
	}
}
