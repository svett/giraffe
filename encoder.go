package giraffe

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// ContentJSON header value for JSON data.
	ContentJSON = "application/json"
	// ContentJSONP header value for JSONP data.
	ContentJSONP = "application/javascript"
	// ContentType header constant.
	ContentType = "Content-Type"
	// ContentDefaultCharset default character encoding.
	ContentDefaultCharset = "UTF-8"
)

// ContentTypeWithCharset returns the contentype with the default charset
func ContentTypeWithCharset(contentType string) string {
	return fmt.Sprintf("%s; charset=%s", contentType, ContentDefaultCharset)
}

// Model represents a encoder data
type Model interface{}

// HTTPEncoder encodes into a different formats
type HTTPEncoder struct {
	writer http.ResponseWriter
}

// NewHTTPEncoder creates a new encoder for concrete writer
func NewHTTPEncoder(writer http.ResponseWriter) *HTTPEncoder {
	return &HTTPEncoder{writer: writer}
}

// EncodeJSON encodes a data as json
func (enc *HTTPEncoder) EncodeJSON(model Model) error {
	enc.writer.Header().Set(ContentType, ContentTypeWithCharset(ContentJSON))

	err := json.NewEncoder(enc.writer).Encode(model)
	if err != nil {
		http.Error(enc.writer, fmt.Sprintf("Unable to encode '%v' as JSON data", model), http.StatusInternalServerError)
	}
	return err
}

// EncodeJSONP encodes a data as jsonp
func (enc *HTTPEncoder) EncodeJSONP(callback string, model Model) error {
	enc.writer.Header().Set(ContentType, ContentTypeWithCharset(ContentJSONP))

	data, _ := json.Marshal(model)
	_, err := fmt.Fprintf(enc.writer, "%s(%s)", callback, string(data))
	if err != nil {
		http.Error(enc.writer, fmt.Sprintf("Unable to encode '%v' as JSON for javascript func %s", model, callback), http.StatusInternalServerError)
	}
	return err
}
