package giraffe_test

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/svett/giraffe"
	"github.com/svett/giraffe/fakes"
)

var _ = Describe("HTMLTemplateRenderer", func() {
	var (
		renderer       *giraffe.HTMLTemplateRenderer
		recorder       *httptest.ResponseRecorder
		responseWriter http.ResponseWriter
		provider       *fakes.FakeHTMLTemplateProvider
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		responseWriter = recorder

		templates := template.New("assets")

		buffer, err := ioutil.ReadFile("assets/home.tmpl")
		Expect(err).NotTo(HaveOccurred())

		tmpl := templates.New("home")
		tmpl.Parse(string(buffer))

		provider = new(fakes.FakeHTMLTemplateProvider)
		provider.ProvideReturns(templates, nil)
	})

	JustBeforeEach(func() {
		renderer = giraffe.NewHTMLTemplateRendererWithProvider(responseWriter, provider)
	})

	It("renders the templates", func() {
		Expect(renderer.Render("home", "Ben")).To(Succeed())
		Expect(provider.ProvideCallCount()).To(Equal(1))
		Expect(recorder.Body.String()).To(Equal("Welcome home, Ben!\n"))
	})

	It("has the corrent content type", func() {
		Expect(renderer.Render("home", "Ben")).To(Succeed())
		Expect(recorder.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"text/html; charset=UTF-8"}))
	})

	It("has the correct status code", func() {
		Expect(renderer.Render("home", "Ben")).To(Succeed())
		Expect(recorder.Code).To(Equal(http.StatusOK))
	})

	Context("when template rendering fails", func() {
		var fakeResponseWriter *fakes.FakeResponseWriter

		BeforeEach(func() {
			fakeResponseWriter = fakes.NewFakeResponseWriter(fakes.FuncWriter(func(_ []byte) (int, error) {
				return -1, fmt.Errorf("Oh no!")
			}))
			responseWriter = fakeResponseWriter
		})

		It("returns the error", func() {
			Expect(renderer.Render("home", "Ben")).To(MatchError("Oh no!"))
		})

		It("has correct status code", func() {
			renderer.Render("home", "Ben")
			Expect(fakeResponseWriter.Code()).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when the template provider fails", func() {
		BeforeEach(func() {
			provider.ProvideReturns(nil, errors.New("oh no!"))
		})

		It("returns the error", func() {
			Expect(renderer.Render("home", "Ben")).To(MatchError("oh no!"))
		})

		It("has correct status code", func() {
			renderer.Render("home", "Ben")
			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})

var _ = Describe("NewHTMLTemplateRenderer", func() {
	It("creates a renderer with default provider", func() {
		provider := new(fakes.FakeHTMLTemplateProvider)
		provider.ProvideReturns(nil, errors.New("oh no!"))
		giraffe.SetHTMLTemplateProvider(provider)
		renderer := giraffe.NewHTMLTemplateRenderer(httptest.NewRecorder())
		renderer.Render("my_template", nil)
		Expect(provider.ProvideCallCount()).To(Equal(1))
	})
})
