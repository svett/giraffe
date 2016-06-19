package fakes

import (
	"io"
	"net/http"
)

// FuncWriter converts a func into io.Writer
type FuncWriter func(data []byte) (int, error)

// Write writes a data
func (f FuncWriter) Write(data []byte) (int, error) {
	return f(data)
}

// FakeResponseWriter is a mockup of http.ResponseWriter
type FakeResponseWriter struct {
	code   int
	header http.Header
	buffer io.Writer
}

// NewFakeResponseWriter creates a new FakeResponseWriter
func NewFakeResponseWriter(buffer io.Writer) *FakeResponseWriter {
	return &FakeResponseWriter{
		header: http.Header{},
		buffer: buffer,
	}
}

// Header returns the http.Header
func (w *FakeResponseWriter) Header() http.Header {
	return w.header
}

// Code returns the status code
func (w *FakeResponseWriter) Code() int {
	return w.code
}

// Write writes a data
func (w *FakeResponseWriter) Write(data []byte) (int, error) {
	return w.buffer.Write(data)
}

// WriteHeader write a header
func (w *FakeResponseWriter) WriteHeader(code int) {
	w.code = code
}
