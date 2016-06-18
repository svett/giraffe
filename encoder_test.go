package giraffe_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/svett/giraffe"
	"github.com/svett/giraffe/mocks"
)

var _ = Describe("HTTPEncoder", func() {
	var (
		encoder        *giraffe.HTTPEncoder
		responseWriter http.ResponseWriter
		recoder        *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		recoder = httptest.NewRecorder()
		responseWriter = recoder
	})

	JustBeforeEach(func() {
		encoder = giraffe.NewHTTPEncoder(responseWriter)
	})

	Describe("JSON", func() {
		var model map[string]string

		BeforeEach(func() {
			model = make(map[string]string)
			model["name"] = "Unknown"
		})

		Describe("EncodeJSON", func() {
			It("encodes a json format", func() {
				Expect(encoder.EncodeJSON(model)).To(Succeed())

				var attrib map[string]string
				Expect(json.Unmarshal(recoder.Body.Bytes(), &attrib)).To(Succeed())
				Expect(attrib).To(Equal(model))
			})

			It("has the corrent content type", func() {
				Expect(encoder.EncodeJSON(model)).To(Succeed())
				Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=UTF-8"}))
			})

			It("has the correct status code", func() {
				Expect(encoder.EncodeJSON(model)).To(Succeed())
				Expect(recoder.Code).To(Equal(http.StatusOK))
			})
		})

		Describe("EncodeJSONP", func() {
			It("encodes a json format", func() {
				Expect(encoder.EncodeJSONP("my_callback", model)).To(Succeed())
				Expect(recoder.Body.String()).To(Equal("my_callback({\"name\":\"Unknown\"})"))
			})

			It("has the corrent content type", func() {
				Expect(encoder.EncodeJSONP("my_callback_func", model)).To(Succeed())
				Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"application/javascript; charset=UTF-8"}))
			})

			It("has the correct status code", func() {
				Expect(encoder.EncodeJSONP("my_callback_func", model)).To(Succeed())
				Expect(recoder.Code).To(Equal(http.StatusOK))
			})
		})

		Context("when encoding fails", func() {
			var fakeResponseWriter *mocks.FakeResponseWriter

			BeforeEach(func() {
				fakeResponseWriter = mocks.NewFakeResponseWriter(mocks.FuncWriter(func(_ []byte) (int, error) {
					return -1, fmt.Errorf("Oh no!")
				}))
				responseWriter = fakeResponseWriter
			})

			Describe("EncodeJSON", func() {
				It("returns the error", func() {
					Expect(encoder.EncodeJSON(model)).To(MatchError("Oh no!"))
				})

				It("has correct status code", func() {
					encoder.EncodeJSON(model)
					Expect(fakeResponseWriter.Code()).To(Equal(http.StatusInternalServerError))
				})
			})

			Describe("EncodeJSONP", func() {
				It("returns the error", func() {
					Expect(encoder.EncodeJSONP("my_callback_func", model)).To(MatchError("Oh no!"))
				})

				It("has correct status code", func() {
					encoder.EncodeJSONP("my_callback_func", model)
					Expect(fakeResponseWriter.Code()).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})
})
