package giraffe

import (
	"fmt"
	"net/http"
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
	// ContentXHTML header value for XHTML data.
	ContentXHTML = "application/xhtml+xml"

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
