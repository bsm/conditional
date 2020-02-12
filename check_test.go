package conditional_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/bsm/conditional"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Check", func() {
	Describe("If-Match", func() {
		It("should require match", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-Match", `"123456"`)

			w := httptest.NewRecorder()
			Expect(conditional.Check(w, r)).To(BeTrue())
			Expect(w.Code).To(Equal(http.StatusPreconditionFailed))

			w = httptest.NewRecorder()
			w.Header().Set("ETag", `"123456"`)
			Expect(conditional.Check(w, r)).To(BeFalse())
		})

		It("should support wildcards", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-Match", `*`)

			w := httptest.NewRecorder()
			Expect(conditional.Check(w, r)).To(BeFalse())
		})

		It("should support multiple values", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-Match", `"a","b","c"`)

			w := httptest.NewRecorder()
			w.Header().Set("ETag", `"b"`)
			Expect(conditional.Check(w, r)).To(BeFalse())
		})
	})

	Describe("If-None-Match", func() {
		It("should require non-match", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-None-Match", `"123456"`)

			w := httptest.NewRecorder()
			Expect(conditional.Check(w, r)).To(BeFalse())

			w = httptest.NewRecorder()
			w.Header().Set("ETag", `"123456"`)
			Expect(conditional.Check(w, r)).To(BeTrue())
			Expect(w.Code).To(Equal(http.StatusNotModified))
		})

		It("should support wildcards", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-None-Match", `*`)

			w := httptest.NewRecorder()
			Expect(conditional.Check(w, r)).To(BeTrue())
		})

		It("should support multiple values", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-None-Match", `"a","b","c"`)

			w := httptest.NewRecorder()
			w.Header().Set("ETag", `"b"`)
			Expect(conditional.Check(w, r)).To(BeTrue())
		})
	})

	Describe("If-Modified-Since", func() {
		It("should check", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-Modified-Since", `Mon, 20 Jan 2020 20:01:20 GMT`)

			w := httptest.NewRecorder()
			Expect(conditional.Check(w, r)).To(BeFalse())

			w = httptest.NewRecorder()
			w.Header().Set("Last-Modified", `Mon, 20 Jan 2020 20:01:20 GMT`)
			Expect(conditional.Check(w, r)).To(BeTrue())
			Expect(w.Code).To(Equal(http.StatusNotModified))
		})

		It("should not trigger when If-None-Match set", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-None-Match", `"123456"`)
			r.Header.Set("If-Modified-Since", `Mon, 20 Jan 2020 20:01:20 GMT`)

			w := httptest.NewRecorder()
			w.Header().Set("ETag", `"123457"`)
			w.Header().Set("Last-Modified", `Mon, 20 Jan 2020 20:01:20 GMT`)
			Expect(conditional.Check(w, r)).To(BeFalse())
		})
	})

	Describe("If-Unmodified-Since", func() {
		It("should check", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-Unmodified-Since", `Mon, 20 Jan 2020 20:01:20 GMT`)

			w := httptest.NewRecorder()
			Expect(conditional.Check(w, r)).To(BeFalse())

			w = httptest.NewRecorder()
			w.Header().Set("Last-Modified", `Mon, 20 Jan 2020 20:01:20 GMT`)
			Expect(conditional.Check(w, r)).To(BeFalse())

			w = httptest.NewRecorder()
			w.Header().Set("Last-Modified", `Mon, 20 Jan 2020 20:01:21 GMT`)
			Expect(conditional.Check(w, r)).To(BeTrue())
			Expect(w.Code).To(Equal(http.StatusPreconditionFailed))
		})

		It("should not trigger when If-Match set", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("If-Match", `"123456"`)
			r.Header.Set("If-Unmodified-Since", `Mon, 20 Jan 2020 20:01:20 GMT`)

			w := httptest.NewRecorder()
			w.Header().Set("ETag", `"123456"`)
			w.Header().Set("Last-Modified", `Mon, 20 Jan 2020 20:01:21 GMT`)
			Expect(conditional.Check(w, r)).To(BeFalse())
		})
	})
})
