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

func main() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set ETag and/or Last-Modified headers
		w.Header().Set("ETag", `"strong"`)
		w.Header().Set("Last-Modified", "Fri, 05 Jan 2018 11:25:15 GMT")

		// perform conditional check
		if conditional.Check(w, r) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := new(http.Client)

	// make a plain GET request
	req, err := http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Status)	// => 204 No Content

	// now, try it with a matchingg "If-Modified-Since"
	req, err = http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("If-Modified-Since", "Fri, 05 Jan 2018 11:25:15 GMT")
	res, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Status)

}
```
