package edit

import (
	"context"
	"io"
	"log"
	
	writer "github.com/whosonfirst/go-writer/v3"
	cli_writer "github.com/whosonfirst/wof/writer"	
)

type EditWriter struct {
	writer.Writer
}

func NewEditWriter(ctx context.Context, uri string) (writer.Writer, error) {

	wr := &EditWriter{}
	return wr, nil
}

// Write copies the content of 'fh' to 'path' using an `io.Discard` writer.
func (wr *EditWriter) Write(ctx context.Context, path string, r io.ReadSeeker) (int64, error) {

	log.Println("WRITE", path)
	
	body, err := io.ReadAll(r)

	if err != nil {
		return 0, err
	}
	
	err = cli_writer.Write(ctx, path, body)

	if err != nil {
		return 0, err
	}

	return int64(len(body)), nil
}

// WriterURI returns the value of 'path'
func (wr *EditWriter) WriterURI(ctx context.Context, path string) string {
	return path
}

// Flush is a no-op to conform to the `Writer` instance and returns nil.
func (wr *EditWriter) Flush(ctx context.Context) error {
	return nil
}

// Close is a no-op to conform to the `Writer` instance and returns nil.
func (wr *EditWriter) Close(ctx context.Context) error {
	return nil
}

// SetLogger is a no-op to conform to the `Writer` instance and returns nil.
func (wr *EditWriter) SetLogger(ctx context.Context, logger *log.Logger) error {
	return nil
}
