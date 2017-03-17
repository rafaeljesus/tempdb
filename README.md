## Tempdb

* TempDB is Redis-backed temporary key-value store for Go.
* Useful for storing temporary data such as login codes, authentication tokens, and temporary passwords.
* A Go version of [tempDB](https://github.com/shanev/tempdb)

## Installation
```bash
go get -u github.com/rafaeljesus/tempdb
```

## Usage
Tempdb stores an expiring (or non-expiring) key/value pair in Redis.

### Tempdb
```go
import "github.com/rafaeljesus/tempdb"

temp, err := tempdb.New(tempdb.Options{
  Addr: "localhost:6379",
  Password: "foo",
})

if err = temp.Insert("key", "value", 0); err != nil {
  // handle failure insert key
}

if err = temp.Insert("key2", "value", time.Hour); err != nil {
  // handle failure insert key
}

if err = temp.Get("key"); err != nil {
  // handle failure to get value
}

```

## Contributing
- Fork it
- Create your feature branch (`git checkout -b my-new-feature`)
- Commit your changes (`git commit -am 'Add some feature'`)
- Push to the branch (`git push origin my-new-feature`)
- Create new Pull Request

## Badges

[![Build Status](https://circleci.com/gh/rafaeljesus/tempdb.svg?style=svg)](https://circleci.com/gh/rafaeljesus/tempdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/rafaeljesus/tempdb)](https://goreportcard.com/report/github.com/rafaeljesus/tempdb)
[![Go Doc](https://godoc.org/github.com/rafaeljesus/tempdb?status.svg)](https://godoc.org/github.com/rafaeljesus/tempdb)

---

> GitHub [@rafaeljesus](https://github.com/rafaeljesus) &nbsp;&middot;&nbsp;
> Medium [@_jesus_rafael](https://medium.com/@_jesus_rafael) &nbsp;&middot;&nbsp;
> Twitter [@_jesus_rafael](https://twitter.com/_jesus_rafael)
