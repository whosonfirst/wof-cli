package writer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/natefinch/atomic"
	wof_writer "github.com/whosonfirst/go-writer/v3"
)

const STDOUT string = "-"

// Write writes the data in 'body' to 'uri'.
func Write(ctx context.Context, uri string, body []byte) error {

	br := bytes.NewReader(body)
	return WriteReader(ctx, uri, br)
}

// Write writes the data in 'r' to 'uri'.
func WriteReader(ctx context.Context, uri string, r io.Reader) error {

	if uri == STDOUT {

		_, err := io.Copy(os.Stdout, r)

		if err != nil {
			return fmt.Errorf("Failed to write body, %w", err)
		}

		return nil
	}

	return atomic.WriteFile(uri, r)
}

// Writer implements the `whosonfirst/go-writer/v3.Writer` interface.
type Writer struct {
	wof_writer.Writer
}

// NewWriter returns a `Writer` instance implementing the `whosonfirst/go-writer/v3.Writer` interface.
func NewWriter() (wof_writer.Writer, error) {
	wr := &Writer{}
	return wr, nil
}

// Write copies the content of 'fh' to 'path' using an `io.Discard` writer.
func (wr *Writer) Write(ctx context.Context, path string, r io.ReadSeeker) (int64, error) {

	err := WriteReader(ctx, path, r)

	if err != nil {
		return 0, err
	}

	return 0, nil
}

// WriterURI returns the value of 'path'
func (wr *Writer) WriterURI(ctx context.Context, path string) string {
	return path
}

// Flush is a no-op to conform to the `Writer` instance and returns nil.
func (wr *Writer) Flush(ctx context.Context) error {
	return nil
}

// Close is a no-op to conform to the `Writer` instance and returns nil.
func (wr *Writer) Close(ctx context.Context) error {
	return nil
}

// SetLogger is a no-op to conform to the `Writer` instance and returns nil.
func (wr *Writer) SetLogger(ctx context.Context, logger *log.Logger) error {
	return nil
}
