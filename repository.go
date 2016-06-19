package giraffe

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TemplateCompilation defines the compilcation options
type TemplateCompilation uint8

const (
	// CompileAlways enables compilation everytime when the templates are used
	CompileAlways TemplateCompilation = iota
	// CompileOnce  enables compilation only once when the templates are used for first time
	CompileOnce
)

// HTMLTemplateRepository represents a template repository
type HTMLTemplateRepository struct {
	// templates are compiled HTML templates
	templates *template.Template

	// Directory to load templates. Default is "templates".
	Directory string
	// FileExtensions to parse template files from. Defaults to [".tmpl"].
	FileExtension string
	// Compilation to compile the templates
	Compilation TemplateCompilation
	// UtilFuncs is a FuncMap to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
	UtilFuncs template.FuncMap
}

// Provide returns the repository compiled templates
func (repository *HTMLTemplateRepository) Provide() (*template.Template, error) {
	option := repository.Compilation
	if (option == CompileOnce && repository.templates == nil) || option == CompileAlways {
		repository.templates = template.New(repository.Directory)

		filepath.Walk(repository.Directory, func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				return nil
			}

			rel, ext, err := ext(repository.Directory, path)
			if err != nil {
				return err
			}

			if ext != repository.FileExtension {
				return nil
			}

			tmpl := repository.templates.New(name(rel, ext))
			buffer, _ := ioutil.ReadFile(path)
			tmpl.Funcs(repository.UtilFuncs).Parse(string(buffer))
			return nil
		})
	}

	return repository.templates, nil
}
