package giraffe_test

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/svett/giraffe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("HTMLTemplateRepository", func() {
	var repository *giraffe.HTMLTemplateRepository

	BeforeEach(func() {
		repository = &giraffe.HTMLTemplateRepository{
			Directory:     "assets",
			FileExtension: ".tmpl",
			Compilation:   giraffe.CompileOnce,
			UtilFuncs: template.FuncMap{
				"say_hello": func() string {
					return "Hello, World!"
				},
			},
		}
	})

	It("compiles all templates for given directory", func() {
		templates, err := repository.Provide()
		Expect(err).NotTo(HaveOccurred())

		homeBuffer := gbytes.NewBuffer()
		Expect(templates.ExecuteTemplate(homeBuffer, "home", "John")).To(Succeed())
		Expect(homeBuffer).To(gbytes.Say("Welcome home, John!"))

		contentBuffer := gbytes.NewBuffer()
		Expect(templates.ExecuteTemplate(contentBuffer, "content", "Bible")).To(Succeed())
		Expect(contentBuffer).To(gbytes.Say("Content of Bible"))
	})

	It("compiles all templates with particular file extension", func() {
		templates, err := repository.Provide()
		Expect(err).NotTo(HaveOccurred())
		Expect(templates.Lookup("assets/info.notmpl")).To(BeNil())
	})

	Context("when the directory has subdirectory with template extension name", func() {
		BeforeEach(func() {
			var err error
			repository.Directory, err = ioutil.TempDir("", "templates")
			Expect(err).NotTo(HaveOccurred())

			subdirectory := filepath.Join(repository.Directory, "index.tmpl")
			Expect(os.Mkdir(subdirectory, 0755)).To(Succeed())

			templatePath := filepath.Join(subdirectory, "page.tmpl")

			Expect(ioutil.WriteFile(templatePath, []byte("Index, {{.}}"), 0777)).To(Succeed())
		})

		It("compiles all templates for given directory", func() {
			templates, err := repository.Provide()
			Expect(err).NotTo(HaveOccurred())

			homeBuffer := gbytes.NewBuffer()
			Expect(templates.ExecuteTemplate(homeBuffer, "index.tmpl/page", "John")).To(Succeed())
			Expect(homeBuffer).To(gbytes.Say("Index, John"))
		})
	})

	It("compiles all templates once", func() {
		var (
			err       error
			templates *template.Template
		)

		repository.Directory, err = ioutil.TempDir("", "templates")
		Expect(err).NotTo(HaveOccurred())

		templatePath := filepath.Join(repository.Directory, "page.tmpl")
		Expect(ioutil.WriteFile(templatePath, []byte("Index, {{.}}"), 0777)).To(Succeed())

		templates, err = repository.Provide()
		Expect(err).NotTo(HaveOccurred())

		buffer := gbytes.NewBuffer()
		Expect(templates.ExecuteTemplate(buffer, "page", "John")).To(Succeed())
		Expect(buffer).To(gbytes.Say("Index, John"))

		Expect(ioutil.WriteFile(templatePath, []byte("Welcome, {{.}}"), 0777)).To(Succeed())
		templates, err = repository.Provide()
		Expect(err).NotTo(HaveOccurred())

		buffer = gbytes.NewBuffer()
		Expect(templates.ExecuteTemplate(buffer, "page", "John")).To(Succeed())
		Expect(buffer).To(gbytes.Say("Index, John"))
	})

	It("provides an utility functions to the template", func() {
		templates, err := repository.Provide()
		Expect(err).NotTo(HaveOccurred())

		homeBuffer := gbytes.NewBuffer()
		Expect(templates.ExecuteTemplate(homeBuffer, "utils", nil)).To(Succeed())
		Expect(homeBuffer).To(gbytes.Say("Hello, World!"))
	})

	Context("when template compilation is set to 'always'", func() {
		It("compiles the templates everytime", func() {
			var (
				err       error
				templates *template.Template
			)

			repository.Compilation = giraffe.CompileAlways
			repository.Directory, err = ioutil.TempDir("", "templates")
			Expect(err).NotTo(HaveOccurred())

			templatePath := filepath.Join(repository.Directory, "page.tmpl")
			Expect(ioutil.WriteFile(templatePath, []byte("Index, {{.}}"), 0777)).To(Succeed())

			templates, err = repository.Provide()
			Expect(err).NotTo(HaveOccurred())

			buffer := gbytes.NewBuffer()
			Expect(templates.ExecuteTemplate(buffer, "page", "John")).To(Succeed())
			Expect(buffer).To(gbytes.Say("Index, John"))

			Expect(ioutil.WriteFile(templatePath, []byte("Welcome, {{.}}"), 0777)).To(Succeed())
			templates, err = repository.Provide()
			Expect(err).NotTo(HaveOccurred())

			buffer = gbytes.NewBuffer()
			Expect(templates.ExecuteTemplate(buffer, "page", "John")).To(Succeed())
			Expect(buffer).To(gbytes.Say("Welcome, John"))
		})
	})
})
