package giraffe

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

var (
	mu              sync.RWMutex
	defaultProvider HTMLTemplateProvider
)

func init() {
	defaultProvider = &HTMLTemplateRepository{
		Directory:     "templates",
		FileExtension: ".tmpl",
		Compilation:   CompileOnce,
	}
}

//go:generate counterfeiter -o fakes/fake_html_template_provider.go . HTMLTemplateProvider

// HTMLTemplateProvider provides a templates
type HTMLTemplateProvider interface {
	Provide() (*template.Template, error)
}

// HTMLTemplateRenderer renders a templates of repository
type HTMLTemplateRenderer struct {
	writer   http.ResponseWriter
	provider HTMLTemplateProvider
}

// Render renders a template
func (renderer *HTMLTemplateRenderer) Render(template string, model Model) error {
	templates, err := renderer.provider.Provide()
	if err != nil {
		renderer.errorf(template, err)
		return err
	}
	setContentType(renderer.writer, ContentHTML)
	err = templates.ExecuteTemplate(renderer.writer, template, model)
	if err != nil {
		renderer.errorf(template, err)
		return err
	}
	return nil
}

func (renderer *HTMLTemplateRenderer) errorf(template string, err error) {
	http.Error(renderer.writer, fmt.Sprintf("Unable to render '%s' html template: %s", template, err.Error()), http.StatusInternalServerError)
}

// NewHTMLTemplateRendererWithProvider create a new HTMLTemplateRenderer for specific provider
func NewHTMLTemplateRendererWithProvider(writer http.ResponseWriter, provider HTMLTemplateProvider) *HTMLTemplateRenderer {
	return &HTMLTemplateRenderer{
		writer:   writer,
		provider: provider,
	}
}

// NewHTMLTemplateRenderer create a new HTMLTemplateRenderer
func NewHTMLTemplateRenderer(writer http.ResponseWriter) *HTMLTemplateRenderer {
	return NewHTMLTemplateRendererWithProvider(writer, defaultProvider)
}

// SetHTMLTemplateProvider sets the default repository
func SetHTMLTemplateProvider(provider HTMLTemplateProvider) {
	mu.Lock()
	defer mu.Unlock()
	defaultProvider = provider
}
