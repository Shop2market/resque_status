package resque_status_test

import (
	. "github.com/Shop2market/resque_status"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lock", func() {
	Context("Lock key", func() {
		It("generates lock keys, sorting params", func() {
			lock := &Lock{LockKeyPrefix: "gen", KeyParamNames: []string{"shop_id", "channel_id"}}
			params := map[string]interface{}{"shop_id": 10, "channel_id": 457, "time_id": "2015", "debug": true}
			Expect(lock.Key(params)).To(Equal("resque:lock:gen-channel_id=457|shop_id=10"))
		})
	})
})
