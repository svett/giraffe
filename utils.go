package giraffe

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	// ContentBinary header value for binary data.
	ContentBinary = "application/octet-stream"
	// ContentJSON header value for JSON data.
	ContentJSON = "application/json"
	// ContentJSONP header value for JSONP data.
	ContentJSONP = "application/javascript"
	// ContentText header value for Text data.
	ContentText = "text/plain"
	// ContentHTML header value for HTML data.
	ContentHTML = "text/html"

	// ContentType header constant.
	ContentType = "Content-Type"
	// ContentDefaultCharset default character encoding.
	ContentDefaultCharset = "UTF-8"
)

func setContentType(writer http.ResponseWriter, contentType string) {
	if writer.Header().Get(ContentType) != "" {
		return
	}

	contentType = fmt.Sprintf("%s; charset=%s", contentType, ContentDefaultCharset)
	writer.Header().Set(ContentType, contentType)
}

func name(dir, ext string) string {
	name := (dir[0 : len(dir)-len(ext)])
	return name
}

func ext(dir, path string) (string, string, error) {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return "", "", err
	}

	ext := ""
	if strings.Index(rel, ".") != -1 {
		ext = filepath.Ext(rel)
	}
	return rel, ext, nil
}
