package resque_status_test

import (
	"flag"

	"github.com/Shop2market/goworker"
	. "github.com/Shop2market/resque_status"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Enqueue", func() {
	var redisConnection *goworker.RedisConn
	BeforeEach(func() {
		flag.Set("exit-on-complete", "true")
		flag.Set("queues", "test_queue")
		flag.Set("concurrency", "1")
		flag.Set("use-number", "true")
		err := goworker.Init()
		if err != nil {
			panic(err)
		}
		redisConnection, err = goworker.GetConn()
		if err != nil {
			panic(err)
		}
		redisConnection.Do("FLUSHDB")
		GenerateUUID = func() string {
			return "NEW_UUID"
		}
	})
	Context("Lock", func() {
		It("creates lock for a job if defined one", func() {
			enq := NewEnqueuer("test", "gen", "Job::Class", []string{"id", "year"})
			params := map[string]interface{}{"id": 10, "year": "2015", "debug": true}
			Expect(enq.Enqueue(params)).NotTo(HaveOccurred())
			Expect(redisConnection.Do("GET", "resque:lock:Job::Class-id=10|year=2015")).To(BeEquivalentTo(`NEW_UUID`))
			Expect(enq.Enqueue(params)).NotTo(HaveOccurred())
		})
	})

})
