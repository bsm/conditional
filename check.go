package conditional

import (
	"net/http"
	"net/textproto"
	"time"
)

// Check evaluates request preconditions and return true when a precondition
// resulted in sending http.StatusNotModified or http.StatusPreconditionFailed.
//
// This method relies on ETag and Last-Modified headers being already set on the
// http.ResponseWriter.
func Check(w http.ResponseWriter, r *http.Request) (done bool) {
	return CheckStatus(w, r) != 0
}

// CheckStatus behaves like Check, but returns the http status code
// that was used for the response instead of a boolean.
//
// 	http.StatusNotModified
// 	http.StatusPreconditionFailed
//	0 = response header was not written
//
func CheckStatus(w http.ResponseWriter, r *http.Request) (code int) {
	code = Evaluate(w.Header(), r)
	switch code {
	case http.StatusNotModified:
		NotModified(w)
	case http.StatusPreconditionFailed:
		w.WriteHeader(http.StatusPreconditionFailed)
	}
	return
}

// Evaluate evaluates request preconditions and return the appropriate
// http status code, which is one of the following:
//
// 	http.StatusNotModified - content was not modified
// 	http.StatusPreconditionFailed - requested content was modified meanwhile
//	0 - otherwise.
//
// Unlike Check it does not automatically write a response to the client.
// Like Check, it relies on ETag and Last-Modified headers being already
// set.
func Evaluate(h http.Header, r *http.Request) int {
	etag := h.Get("ETag")
	modTime, _ := http.ParseTime(h.Get("Last-Modified"))

	// This function carefully follows RFC 7232 section 6.
	ch := checkIfMatch(r, etag)
	if ch == condNone {
		ch = checkIfUnmodifiedSince(r, modTime)
	}
	if ch == condFalse {
		return http.StatusPreconditionFailed
	}

	switch checkIfNoneMatch(r, etag) {
	case condFalse:
		if r.Method == "GET" || r.Method == "HEAD" {
			return http.StatusNotModified
		}
		return http.StatusPreconditionFailed
	case condNone:
		if checkIfModifiedSince(r, modTime) == condFalse {
			return http.StatusNotModified
		}
	}

	return 0
}

// condResult is the result of an HTTP request precondition check.
// See https://tools.ietf.org/html/rfc7232 section 3.
type condResult uint8

const (
	condNone condResult = iota
	condTrue
	condFalse
)

func checkIfMatch(r *http.Request, etag string) condResult {
	val := r.Header.Get("If-Match")
	if val == "" {
		return condNone
	}

	for {
		val = textproto.TrimString(val)
		if len(val) == 0 {
			break
		} else if val[0] == ',' {
			val = val[1:]
			continue
		} else if val[0] == '*' {
			return condTrue
		}

		t, remain := ScanETag(val)
		if t == "" {
			break
		} else if t.IsStrongMatch(etag) {
			return condTrue
		}
		val = remain
	}
	return condFalse
}

func checkIfNoneMatch(r *http.Request, etag string) condResult {
	val := r.Header.Get("If-None-Match")
	if val == "" {
		return condNone
	}

	for {
		val = textproto.TrimString(val)
		if len(val) == 0 {
			break
		} else if val[0] == ',' {
			val = val[1:]
		} else if val[0] == '*' {
			return condFalse
		}

		t, remain := ScanETag(val)
		if t == "" {
			break
		} else if t.IsWeakMatch(etag) {
			return condFalse
		}
		val = remain
	}
	return condTrue
}

func checkIfUnmodifiedSince(r *http.Request, modTime time.Time) condResult {
	val := r.Header.Get("If-Unmodified-Since")
	if val == "" || isZeroTime(modTime) {
		return condNone
	}

	if t, err := http.ParseTime(val); err == nil {
		// The Date-Modified header truncates sub-second precision, so
		// use mtime < t+1s instead of mtime <= t to check for unmodified.
		if modTime.Before(t.Add(1 * time.Second)) {
			return condTrue
		}
		return condFalse
	}
	return condNone
}

func checkIfModifiedSince(r *http.Request, modTime time.Time) condResult {
	if r.Method != "GET" && r.Method != "HEAD" {
		return condNone
	}

	val := r.Header.Get("If-Modified-Since")
	if val == "" || isZeroTime(modTime) {
		return condNone
	}

	t, err := http.ParseTime(val)
	if err != nil {
		return condNone
	}

	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if modTime.Before(t.Add(1 * time.Second)) {
		return condFalse
	}
	return condTrue
}
