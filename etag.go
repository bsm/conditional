package conditional

import (
	"net/textproto"
	"strings"
)

// ETag represents an ETag string.
type ETag string

// ScanETag determines if a syntactically valid ETag is present at s. If so,
// the ETag and remaining text after consuming ETag is returned. Otherwise,
// it returns "", "".
//
// Taken from https://golang.org/src/net/http/fs.go.
// Copyright 2009 The Go Authors. All rights reserved.
func ScanETag(s string) (etag ETag, remain string) {
	s = textproto.TrimString(s)
	start := 0
	if strings.HasPrefix(s, "W/") {
		start = 2
	}
	if len(s[start:]) < 2 || s[start] != '"' {
		return "", ""
	}
	// ETag is either W/"text" or "text".
	// See RFC 7232 2.3.
	for i := start + 1; i < len(s); i++ {
		c := s[i]
		switch {
		// Character values allowed in ETags.
		case c == 0x21 || c >= 0x23 && c <= 0x7E || c >= 0x80:
		case c == '"':
			return ETag(s[:i+1]), s[i+1:]
		default:
			return "", ""
		}
	}
	return "", ""
}

// IsStrongMatch returns true if tag matches other strong tag.
func (t ETag) IsStrongMatch(other string) bool {
	return string(t) == other && t != "" && t[0] == '"'
}

// IsWeakMatch returns true if tag matches other tag.
func (t ETag) IsWeakMatch(other string) bool {
	return strings.TrimPrefix(string(t), "W/") == strings.TrimPrefix(other, "W/")
}

// String returns the ETag as string.
func (t ETag) String() string {
	return string(t)
}
