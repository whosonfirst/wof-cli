package writer

import (
	"context"
	"os"
	"fmt"
)

func Write(ctx context.Context, uri string, body []byte, overwrite bool) error {

	_, err := os.Stdout.Write(body)

	if err != nil {
		return fmt.Errorf("Failed to write body, %w", err)
	}

	return nil
}
