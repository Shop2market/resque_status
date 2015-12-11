package resque_status_test

import (
	"encoding/json"

	. "github.com/Shop2market/resque_status"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParamsParser", func() {
	Context("Parse Int", func() {
		It("parses int from params", func() {
			params := map[string]interface{}{"shop_id": json.Number("10")}
			actualParsed, err := ParseInt(params, "shop_id")
			Expect(err).NotTo(HaveOccurred())
			Expect(actualParsed).To(Equal(10))
		})
	})
	Context("Parse Int array", func() {
		It("parses int from params", func() {
			params := map[string]interface{}{"channel_ids": []interface{}{json.Number("10"), json.Number("20")}}
			actualParsed, err := ParseIntArray(params, "channel_ids")
			Expect(err).NotTo(HaveOccurred())
			Expect(actualParsed).To(Equal([]int{10, 20}))
		})
	})

})
