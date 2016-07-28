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
	Context("Parse Json array", func() {
		It("parses misc data types", func() {
			params := map[string]interface{}{"json_data": `[{"W":"1", "P": 100.10, "B":true},{"Winkelproductcode":"12346", "Titel":"WW","PriceIn": 200}]`}
			actualParsed, err := ParseJsonParam(params, "json_data")
			Expect(err).NotTo(HaveOccurred())
			Expect(actualParsed).To(Equal([]map[string]string{
				map[string]string{"W": "1", "P": "100.1", "B": "true"},
				map[string]string{"Winkelproductcode": "12346", "Titel": "WW", "PriceIn": "200"},
			}))
		})
	})
})
