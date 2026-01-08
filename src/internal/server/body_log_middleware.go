package server

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/hlog"
)

type bodyCaptureWriter struct {
	http.ResponseWriter
	status    int
	buf       bytes.Buffer
	maxBytes  int
	truncated bool
}

func (w *bodyCaptureWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *bodyCaptureWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	if w.maxBytes > 0 && w.buf.Len() < w.maxBytes {
		remaining := w.maxBytes - w.buf.Len()
		if remaining < len(p) {
			w.buf.Write(p[:remaining])
			w.truncated = true
		} else {
			w.buf.Write(p)
		}
	} else if w.maxBytes > 0 {
		w.truncated = true
	}

	return w.ResponseWriter.Write(p)
}

func bodyLogMiddleware(enabled bool, maxBytes int) func(http.Handler) http.Handler {
	if !enabled {
		return func(next http.Handler) http.Handler { return next }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqBody, reqTruncated := readRequestBody(r, maxBytes)

			wrapped := &bodyCaptureWriter{
				ResponseWriter: w,
				maxBytes:       maxBytes,
			}

			next.ServeHTTP(wrapped, r)

			logger := hlog.FromRequest(r)

			if isLoggableContentType(r.Header.Get("Content-Type")) && len(reqBody) > 0 {
				logger.Debug().
					Int("body_bytes", len(reqBody)).
					Bool("body_truncated", reqTruncated).
					Str("request_body", string(reqBody)).
					Msg("http request body")
			}

			respStatus := wrapped.status
			if respStatus == 0 {
				respStatus = http.StatusOK
			}
			respBody := wrapped.buf.Bytes()
			if isLoggableContentType(w.Header().Get("Content-Type")) && len(respBody) > 0 {
				logger.Debug().
					Int("status", respStatus).
					Int("body_bytes", len(respBody)).
					Bool("body_truncated", wrapped.truncated).
					Str("response_body", string(respBody)).
					Msg("http response body")
			}
		})
	}
}

func readRequestBody(r *http.Request, maxBytes int) ([]byte, bool) {
	if r.Body == nil {
		return nil, false
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		r.Body = io.NopCloser(bytes.NewBuffer(nil))
		return nil, false
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if maxBytes <= 0 || len(body) <= maxBytes {
		return body, false
	}

	return body[:maxBytes], true
}

func isLoggableContentType(contentType string) bool {
	ct := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	if ct == "" {
		return true
	}
	if strings.HasPrefix(ct, "text/") {
		return true
	}
	switch ct {
	case "application/json",
		"application/xml",
		"text/plain",
		"text/csv",
		"application/x-www-form-urlencoded":
		return true
	default:
		return false
	}
}
