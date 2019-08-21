[![CircleCI](https://circleci.com/gh/spatialcurrent/go-counter/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/go-counter/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/go-counter)](https://goreportcard.com/report/spatialcurrent/go-counter)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/go-counter?status.svg)](https://godoc.org/github.com/spatialcurrent/go-counter) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/go-counter/blob/master/LICENSE.md)

# go-counter

# Description

**go-counter** is a command line program and package for generating frequency distributions.  The **gocounter** command line program supports the following operating systems and architectures.

| GOOS | GOARCH |
| ---- | ------ |
| darwin | amd64 |
| linux | amd64 |
| windows | amd64 |
| linux | arm64 |

# Installation

No installation is required.  Just grab a [release](https://github.com/spatialcurrent/go-counter/releases).  You might want to rename your binary to just `gocounter` (or `counter`) for convenience.

If you do have go already installed, you can just run using `go run main.go` or install with `make install`.

# Usage

### Go

You can import **go-counter** as a library with:

```go
import (
  "github.com/spatialcurrent/go-counter/pkg/counter"
)
```

See [counter](https://godoc.org/github.com/spatialcurrent/go-counter/pkg/counter) in GoDoc for information on how to use Go API.

### CLI

On the command line use `gocounter --help` to view usage.

# Examples

### Go

See the examples in the [counter](https://godoc.org/github.com/spatialcurrent/go-counter/pkg/counter) package documentation.

### CLI

To print the 10 most frequent lines in a file as a JSON array use:

```shell
gocounter top -sejl -n 10 path/to/file
```

# Building

You can build all the released artifacts using `make build` or run the make target for a specific operating system and architecture.

# Testing

To run Go tests use `make test` (or `bash scripts/test.sh`), which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [ineffassign](https://github.com/gordonklaus/ineffassign), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-counter/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
