package conditional_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/bsm/conditional"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotModified", func() {
	It("should clear headers and respond", func() {
		w := httptest.NewRecorder()
		w.HeaderMap.Set("Content-Type", "text/plain")
		w.HeaderMap.Set("Content-Length", "33")
		w.HeaderMap.Set("ETag", `"strong"`)
		w.HeaderMap.Set("Last-Modified", "Fri, 05 Jan 2018 11:25:15 GMT")

		conditional.NotModified(w)
		Expect(w.Code).To(Equal(http.StatusNotModified))
		Expect(w.HeaderMap).To(Equal(http.Header{
			"Etag": {`"strong"`},
		}))
	})

	It("should clear last-modified only if ETag is set", func() {
		w := httptest.NewRecorder()
		w.HeaderMap.Set("Content-Type", "text/plain")
		w.HeaderMap.Set("Content-Length", "33")
		w.HeaderMap.Set("Last-Modified", "Fri, 05 Jan 2018 11:25:15 GMT")

		conditional.NotModified(w)
		Expect(w.Code).To(Equal(http.StatusNotModified))
		Expect(w.HeaderMap).To(Equal(http.Header{
			"Last-Modified": {"Fri, 05 Jan 2018 11:25:15 GMT"},
		}))
	})
})
