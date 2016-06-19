package giraffe_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/svett/giraffe"
	"github.com/svett/giraffe/fakes"
)

var _ = Describe("HTTPLogger", func() {
	var (
		logHandler giraffe.HandlerFunc
		logger     *fakes.FakeLogger
		request    *http.Request
		writer     http.ResponseWriter
	)

	BeforeEach(func() {
		logger = new(fakes.FakeLogger)
		logHandler = giraffe.NewHTTPLogger(logger)
		Expect(logHandler).NotTo(BeNil())

		writer = httptest.NewRecorder()

		var err error
		request, err = http.NewRequest("GET", "http://example.com/foo", nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("processes the next request", func() {
		processedCnt := 0

		logHandler(writer, request, func(w http.ResponseWriter, req *http.Request) {
			Expect(w).NotTo(BeNil())
			Expect(req).To(Equal(request))
			processedCnt++
		})

		Expect(processedCnt).To(Equal(1))
	})

	It("writes an info message to logrus with the right format", func() {
		logHandler(writer, request, func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		Expect(logger.InfoCallCount()).To(Equal(1))
		msg := logger.InfoArgsForCall(0)

		Expect(msg).To(ContainSubstring(time.Now().Format("2006/01/02")))
		Expect(msg).To(ContainSubstring("200"))
		Expect(msg).To(ContainSubstring("GET"))
		Expect(msg).To(ContainSubstring("/foo"))
	})
})
