package giraffe_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGiraffe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Giraffe Suite")
}
