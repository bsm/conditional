package conditional_test

import (
	"github.com/bsm/conditional"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ETag", func() {
	It("should reject scanning bad etags", func() {
		etag, rest := conditional.ScanETag("bad")
		Expect(etag).To(BeEmpty())
		Expect(rest).To(BeEmpty())
	})

	It("should support scanning strong tags", func() {
		etag, rest := conditional.ScanETag(`"strong"`)
		Expect(etag).To(Equal(conditional.ETag(`"strong"`)))
		Expect(rest).To(BeEmpty())
	})

	It("should support scanning weak tags", func() {
		etag, rest := conditional.ScanETag(`W/"weak"`)
		Expect(etag).To(Equal(conditional.ETag(`W/"weak"`)))
		Expect(rest).To(BeEmpty())
	})

	It("should support scanning multiple tags", func() {
		etag, rest := conditional.ScanETag(`"et1" "et2"`)
		Expect(etag).To(Equal(conditional.ETag(`"et1"`)))
		Expect(rest).To(Equal(` "et2"`))

		etag, rest = conditional.ScanETag(rest)
		Expect(etag).To(Equal(conditional.ETag(`"et2"`)))
		Expect(rest).To(BeEmpty())
	})

	It("should check if strong match", func() {
		Expect(conditional.ETag(`"strong"`).IsStrongMatch(`"strong"`)).To(BeTrue())

		Expect(conditional.ETag(``).IsStrongMatch(``)).To(BeFalse())
		Expect(conditional.ETag(`W/"weak"`).IsStrongMatch(`W/"weak"`)).To(BeFalse())
		Expect(conditional.ETag(`W/"weak"`).IsStrongMatch(`"weak"`)).To(BeFalse())
	})

	It("should check if weak match", func() {
		Expect(conditional.ETag(`"strong"`).IsWeakMatch(`"strong"`)).To(BeTrue())
		Expect(conditional.ETag(``).IsWeakMatch(``)).To(BeTrue())
		Expect(conditional.ETag(`W/"weak"`).IsWeakMatch(`W/"weak"`)).To(BeTrue())
		Expect(conditional.ETag(`W/"weak"`).IsWeakMatch(`"weak"`)).To(BeTrue())

		Expect(conditional.ETag(`W/"weak"`).IsWeakMatch(`W/"other"`)).To(BeFalse())
		Expect(conditional.ETag(`"strong"`).IsWeakMatch(`"other"`)).To(BeFalse())
	})
})
