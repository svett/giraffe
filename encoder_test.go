package giraffe_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/svett/giraffe"
	"github.com/svett/giraffe/fakes"
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

	Describe("EncodeJSON", func() {
		It("encodes a json format", func() {
			model := map[string]string{"name": "root"}
			Expect(encoder.EncodeJSON(model)).To(Succeed())

			var attrib map[string]string
			Expect(json.Unmarshal(recoder.Body.Bytes(), &attrib)).To(Succeed())
			Expect(attrib).To(Equal(model))
		})

		It("has the corrent content type", func() {
			model := map[string]string{"name": "root"}
			Expect(encoder.EncodeJSON(model)).To(Succeed())
			Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=UTF-8"}))
		})

		Context("when the content type is already set", func() {
			It("does not change it", func() {
				model := map[string]string{"name": "root"}
				recoder.Header().Set("Content-Type", "Unknown")
				Expect(encoder.EncodeJSON(model)).To(Succeed())
				Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"Unknown"}))
			})
		})

		It("has the correct status code", func() {
			model := map[string]string{"name": "root"}
			Expect(encoder.EncodeJSON(model)).To(Succeed())
			Expect(recoder.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("EncodeJSONP", func() {
		It("encodes a json format", func() {
			model := map[string]string{"name": "Unknown"}
			Expect(encoder.EncodeJSONP("my_callback", model)).To(Succeed())
			Expect(recoder.Body.String()).To(Equal("my_callback({\"name\":\"Unknown\"})"))
		})

		It("has the corrent content type", func() {
			model := map[string]string{"name": "root"}
			Expect(encoder.EncodeJSONP("my_callback_func", model)).To(Succeed())
			Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"application/javascript; charset=UTF-8"}))
		})

		Context("when the content type is already set", func() {
			It("does not change it", func() {
				model := map[string]string{"name": "root"}
				recoder.Header().Set("Content-Type", "Unknown")
				Expect(encoder.EncodeJSONP("my_callback_func", model)).To(Succeed())
				Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"Unknown"}))
			})
		})

		It("has the correct status code", func() {
			model := map[string]string{"name": "root"}
			Expect(encoder.EncodeJSONP("my_callback_func", model)).To(Succeed())
			Expect(recoder.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("EncodeData", func() {
		It("encodes a binary format", func() {
			Expect(encoder.EncodeData([]byte("hello"))).To(Succeed())
			Expect(recoder.Body.String()).To(Equal("hello"))
		})

		It("has the corrent content type", func() {
			Expect(encoder.EncodeData([]byte("hello"))).To(Succeed())
			Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"application/octet-stream; charset=UTF-8"}))
		})

		Context("when the content type is already set", func() {
			It("does not change it", func() {
				recoder.Header().Set("Content-Type", "Unknown")
				Expect(encoder.EncodeData([]byte("hello"))).To(Succeed())
				Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"Unknown"}))
			})
		})

		It("has the correct status code", func() {
			Expect(encoder.EncodeData([]byte("hello"))).To(Succeed())
			Expect(recoder.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("EncodeText", func() {
		It("encodes a text", func() {
			Expect(encoder.EncodeText("hello")).To(Succeed())
			Expect(recoder.Body.String()).To(Equal("hello"))
		})

		It("has the corrent content type", func() {
			Expect(encoder.EncodeText("hello")).To(Succeed())
			Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"text/plain; charset=UTF-8"}))
		})

		Context("when the content type is already set", func() {
			It("does not change it", func() {
				recoder.Header().Set("Content-Type", "Unknown")
				Expect(encoder.EncodeText("hello")).To(Succeed())
				Expect(recoder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"Unknown"}))
			})
		})

		It("has the correct status code", func() {
			Expect(encoder.EncodeText("hello")).To(Succeed())
			Expect(recoder.Code).To(Equal(http.StatusOK))
		})
	})

	Context("when encoding fails", func() {
		var fakeResponseWriter *fakes.FakeResponseWriter

		BeforeEach(func() {
			fakeResponseWriter = fakes.NewFakeResponseWriter(fakes.FuncWriter(func(_ []byte) (int, error) {
				return -1, fmt.Errorf("Oh no!")
			}))
			responseWriter = fakeResponseWriter
		})

		Describe("EncodeJSON", func() {
			It("returns the error", func() {
				model := map[string]string{"name": "root"}
				Expect(encoder.EncodeJSON(model)).To(MatchError("Oh no!"))
			})

			It("has correct status code", func() {
				model := map[string]string{"name": "root"}
				encoder.EncodeJSON(model)
				Expect(fakeResponseWriter.Code()).To(Equal(http.StatusInternalServerError))
			})
		})

		Describe("EncodeJSONP", func() {
			It("returns the error", func() {
				model := map[string]string{"name": "root"}
				Expect(encoder.EncodeJSONP("my_callback_func", model)).To(MatchError("Oh no!"))
			})

			It("has correct status code", func() {
				model := map[string]string{"name": "root"}
				encoder.EncodeJSONP("my_callback_func", model)
				Expect(fakeResponseWriter.Code()).To(Equal(http.StatusInternalServerError))
			})
		})

		Describe("EncodeData", func() {
			It("returns the error", func() {
				Expect(encoder.EncodeData([]byte("hello"))).To(MatchError("Oh no!"))
			})

			It("has correct status code", func() {
				Expect(encoder.EncodeData([]byte("hello"))).To(MatchError("Oh no!"))
				Expect(fakeResponseWriter.Code()).To(Equal(http.StatusInternalServerError))
			})
		})

		Describe("EncodeText", func() {
			It("returns the error", func() {
				Expect(encoder.EncodeText("hello")).To(MatchError("Oh no!"))
			})

			It("has correct status code", func() {
				Expect(encoder.EncodeText("hello")).To(MatchError("Oh no!"))
				Expect(fakeResponseWriter.Code()).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
