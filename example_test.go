package conditional_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/bsm/conditional"
)

func ExampleCheck_ifModifiedSince() {
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
	fmt.Println(res.Status) // => 204 No Content

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
	fmt.Println(res.Status) // => 304 Not Modified

	// Output:
	// 204 No Content
	// 304 Not Modified
}

func ExampleCheck_ifNoneMatch() {
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
	fmt.Println(res.Status) // => 204 No Content

	// now, try it with a matching "If-None-Match"
	req, err = http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("If-None-Match", `"strong"`)
	res, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Status) // => 304 Not Modified

	// Output:
	// 204 No Content
	// 304 Not Modified
}

func ExampleCheck_ifMatch() {
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
	fmt.Println(res.Status) // => 204 No Content

	// now, try it with a matching "If-Match"
	req, err = http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("If-Match", `"strong"`)
	res, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Status) // => 204 No Content

	// finally, try it with a non-matching "If-Match"
	req, err = http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("If-Match", `"OTHER-TAG"`)
	res, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Status) // => 412 Precondition Failed

	// Output:
	// 204 No Content
	// 204 No Content
	// 412 Precondition Failed
}
