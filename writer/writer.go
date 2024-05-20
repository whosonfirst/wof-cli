package writer

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/natefinch/atomic"
)

const STDOUT string = "-"

func Write(ctx context.Context, uri string, body []byte) error {

	if uri == STDOUT {

		_, err := os.Stdout.Write(body)

		if err != nil {
			return fmt.Errorf("Failed to write body, %w", err)
		}

		return nil
	}

	br := bytes.NewReader(body)
	return atomic.WriteFile(uri, br)
}
