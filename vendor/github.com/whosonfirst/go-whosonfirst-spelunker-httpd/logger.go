package httpd

import (
	"io"
	"log/slog"
	"net/http"
	"os"
)

func LoggerWithRequest(req *http.Request, wr io.Writer) *slog.Logger {

	if wr == nil {
		wr = os.Stderr
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger = logger.With("method", req.Method)
	// logger = logger.With("user agent", req.Header.Get("User-Agent"))
	// logger = logger.With("accept", req.Header.Get("Accept"))
	logger = logger.With("path", req.URL.Path)
	logger = logger.With("remote addr", req.RemoteAddr)
	logger = logger.With("user ip", ReadUserIP(req))

	return logger
}

func ReadUserIP(req *http.Request) string {

	addr := req.Header.Get("X-Real-Ip")

	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
	}

	if addr == "" {
		addr = req.RemoteAddr
	}

	return addr
}
