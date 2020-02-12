# Conditional

[![Build Status](https://travis-ci.org/bsm/conditional.png?branch=master)](https://travis-ci.org/bsm/conditional)
[![GoDoc](https://godoc.org/github.com/bsm/conditional?status.png)](http://godoc.org/github.com/bsm/conditional)
[![Go Report Card](https://goreportcard.com/badge/github.com/bsm/conditional)](https://goreportcard.com/report/github.com/bsm/conditional)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Conditional HTTP helpers for [Go](https://golang.org) stdlib's [net/http](https://golang.org/pkg/net/http) package.

Supported headers:

* `If-Match`
* `If-None-Match`
* `If-Modified-Since`
* `If-Unmodified-Since`

## Example:

```go
import(
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/bsm/conditional"
)

func main() {{ "ExampleCheck_ifModifiedSince" | code }}
```
