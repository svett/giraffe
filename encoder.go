package giraffe

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Model represents a encoder data
type Model interface{}

// HTTPEncoder encodes into a different formats
type HTTPEncoder struct {
	writer http.ResponseWriter
}

// EncodeJSON encodes a data as json
func (enc *HTTPEncoder) EncodeJSON(model Model) error {
	setContentType(enc.writer, ContentJSON)

	err := json.NewEncoder(enc.writer).Encode(model)
	if err != nil {
		http.Error(enc.writer, fmt.Sprintf("Unable to encode '%v' as JSON data: %s", model, err.Error()), http.StatusInternalServerError)
	}
	return err
}

// EncodeJSONP encodes a data as jsonp
func (enc *HTTPEncoder) EncodeJSONP(callback string, model Model) error {
	setContentType(enc.writer, ContentJSONP)

	data, _ := json.Marshal(model)
	_, err := fmt.Fprintf(enc.writer, "%s(%s)", callback, string(data))
	if err != nil {
		http.Error(enc.writer, fmt.Sprintf("Unable to encode '%v' as JSON for javascript func %s: %s", model, callback, err.Error()), http.StatusInternalServerError)
	}
	return err
}

// EncodeData encodes an array of bytes
func (enc *HTTPEncoder) EncodeData(data []byte) error {
	setContentType(enc.writer, ContentBinary)

	_, err := enc.writer.Write(data)
	if err != nil {
		http.Error(enc.writer, fmt.Sprintf("Unable to encode binary data: %s", err.Error()), http.StatusInternalServerError)
	}
	return err
}

// EncodeText encodes a plain text
func (enc *HTTPEncoder) EncodeText(text string) error {
	setContentType(enc.writer, ContentText)

	_, err := fmt.Fprint(enc.writer, text)
	if err != nil {
		http.Error(enc.writer, fmt.Sprintf("Unable to encode text '%s': %s", text, err.Error()), http.StatusInternalServerError)
	}
	return err
}

// NewHTTPEncoder creates a new encoder for concrete writer
func NewHTTPEncoder(writer http.ResponseWriter) *HTTPEncoder {
	return &HTTPEncoder{writer: writer}
}
