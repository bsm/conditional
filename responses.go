package conditional

import "net/http"

// NotModified writes a http.StatusNotModified response to the user following
// instructions from RFC 7232 section 4.1:
//   a sender SHOULD NOT generate representation metadata other than the
//   above listed fields unless said metadata exists for the purpose of
//   guiding cache updates (e.g., Last-Modified might be useful if the
//   response does not have an ETag field).
func NotModified(w http.ResponseWriter) {
	h := w.Header()
	h.Del("Content-Type")
	h.Del("Content-Length")
	if h.Get("ETag") != "" {
		h.Del("Last-Modified")
	}

	w.WriteHeader(http.StatusNotModified)
}
