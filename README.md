# Elapsed time

[![GoDoc](https://godoc.org/github.com/rvflash/elapsed?status.svg)](https://godoc.org/github.com/rvflash/elapsed)
[![Build Status](https://github.com/rvflash/elapsed/workflows/build/badge.svg)](https://github.com/rvflash/elapsed/actions?workflow=build)
[![Code Coverage](https://codecov.io/gh/rvflash/elapsed/branch/master/graph/badge.svg)](https://codecov.io/gh/rvflash/elapsed)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/elapsed?)](https://goreportcard.com/report/github.com/rvflash/elapsed)


Golang package to return the elapsed time since a given time in a human-readable format.


### Installation

```bash
$ go get -u github.com/rvflash/elapsed
```

### Usage

```go
t := time.Now().Add(-time.Hour)
fmt.Println(elapsed.Time(t))
// Output: 1 hour ago

t = time.Now().Add(-time.Hour * 24 * 3)
fmt.Println(elapsed.Time(t))
// Output:  3 days ago

t, _ = time.Parse("2006-02-01", "2049-08-19")
fmt.Println(elapsed.Time(t))
// Output: not yet

t = time.Now().Add(-time.Hour * 24 * 3)
fmt.Println(elapsed.LocalTime(t, "fr"))
// Output: il y a 3 jours
```
