package webhelpers

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

// LogrusLogger is a middleware that will log each request recieved, along with
// some useful information, to the given logger.
func LogrusLogger(logger *logrus.Logger, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		entry := logger.WithFields(logrus.Fields{
			"request": r.RequestURI,
			"method":  r.Method,
			"remote":  r.RemoteAddr,
		})

		if id := r.Header.Get(RequestIDKey); id != "" {
			entry = entry.WithField("request_id", id)
		}

		// Wrap the writer so we can track data information.
		neww := WrapWriter(w)

		// Dispatch to the underlying handler.
		entry.Info("started handling request")
		h.ServeHTTP(neww, r)

		// Log final information.
		entry.WithFields(logrus.Fields{
			"bytes_written": neww.BytesWritten(),
			"status":        neww.Status(),
			"text_status":   http.StatusText(neww.Status()),
			"took":          time.Since(start),
		}).Info("completed handling request")
	}

	return http.HandlerFunc(fn)
}
