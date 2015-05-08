package resque_status_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestResqueStatus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ResqueStatus Suite")
}
