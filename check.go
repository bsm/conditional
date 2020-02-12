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
	etag := w.Header().Get("ETag")
	modTime, _ := http.ParseTime(w.Header().Get("Last-Modified"))

	// This function carefully follows RFC 7232 section 6.
	ch := checkIfMatch(r, etag)
	if ch == condNone {
		ch = checkIfUnmodifiedSince(r, modTime)
	}
	if ch == condFalse {
		w.WriteHeader(http.StatusPreconditionFailed)
		return true
	}

	switch checkIfNoneMatch(r, etag) {
	case condFalse:
		if r.Method == "GET" || r.Method == "HEAD" {
			NotModified(w)
		} else {
			w.WriteHeader(http.StatusPreconditionFailed)
		}
		return true
	case condNone:
		if checkIfModifiedSince(r, modTime) == condFalse {
			NotModified(w)
			return true
		}
	}

	return false
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

var unixEpochTime = time.Unix(0, 0)

// isZeroTime reports whether t is obviously unspecified (either zero or Unix()=0).
func isZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(unixEpochTime)
}
