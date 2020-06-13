// Package logging provids a sensiblely configured logger and an htttp access
// logging Handler.
package logging

import (
	"expvar"
	"net/http"
	"strings"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/rs/zerolog"
)

var (
	reqCounter *expvar.Map                 // nolint:gochecknoglobals
	rps        *ratecounter.RateCounter    // nolint:gochecknoglobals
	bps        *ratecounter.RateCounter    // nolint:gochecknoglobals
	durAvg     *ratecounter.AvgRateCounter // nolint:gochecknoglobals

	bytesSecond = expvar.NewInt(bytesRate)   // nolint:gochecknoglobals
	reqSecond   = expvar.NewInt(requestRate) // nolint:gochecknoglobals
	avgDuration = expvar.NewFloat(avgDur)    // nolint:gochecknoglobals

)

func init() { // nolint:gochecknoinits
	reqCounter = expvar.NewMap("reqCounter")
	rps = ratecounter.NewRateCounter(time.Second)
	bps = ratecounter.NewRateCounter(time.Second)
	durAvg = ratecounter.NewAvgRateCounter(time.Minute)
}

const (
	// TimeFormat is the time format for logging.
	TimeFormat = time.RFC3339Nano

	avgDur      = "duration_1m_avg"
	bytesRate   = "bytes_second"
	requestRate = "requests_second"

	err400    = "400"
	err404    = "404"
	err500    = "500"
	requests  = "requests"
	completed = "completed"
	bytes     = "bytes_transferred"
)

// byteCounter implements an io.Writer wrapping the http.ResponseWriter to
// count bytes written in the response.
type byteCounter struct {
	http.ResponseWriter

	status        int
	responseBytes int64
}

// Write writes bytes to the response writer.
func (ar *byteCounter) Write(b []byte) (int, error) {
	written, err := ar.ResponseWriter.Write(b)
	ar.responseBytes += int64(written)

	return written, err
}

// WriteHeader sets the response status.
func (ar *byteCounter) WriteHeader(status int) {
	ar.status = status
	ar.ResponseWriter.WriteHeader(status)
}

// NewAccessLogger returns a constructed AccessLogger pointer.
func NewAccessLogger(handler http.Handler, logger zerolog.Logger) http.Handler {
	logger = logger.With().Str("type", "access").Logger()

	return &AccessLogger{
		handler: handler,
		logger:  logger,
	}
}

// AccessLogger writes an NCSA combined-ish log record. Note this skips the
// rfc931 ident field.
type AccessLogger struct {
	handler http.Handler
	logger  zerolog.Logger
}

// ServeHTTP makes our type a http.HandlerFunc.
func (al *AccessLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCounter.Add(requests, 1)

	clientIP := r.RemoteAddr

	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	bc := &byteCounter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}

	start := time.Now()

	al.handler.ServeHTTP(bc, r)

	dur := time.Since(start)
	username, _, _ := r.BasicAuth()

	al.logger.Info().
		Str("request_id", r.Header.Get("X-Request-ID")).
		Str("client_ip", clientIP).
		Strs("x_forwarded_for", strings.Split(r.Header.Get("X-Forwarded-For"), ", ")).
		Dur("duration", dur).
		Str("domain", r.Host).
		Str("method", r.Method).
		Str("request_uri", r.RequestURI).
		Str("protocol", r.Proto).
		Int("status", bc.status).
		Int64("response_bytes", bc.responseBytes).
		Str("referrer", r.Referer()).
		Str("user", username).
		Str("user_agent", r.UserAgent()).Msg("")

	switch {
	case bc.status == http.StatusNotFound:
		reqCounter.Add(err404, 1)
	case bc.status >= http.StatusBadRequest && bc.status < http.StatusInternalServerError:
		reqCounter.Add(err400, 1)
	case bc.status >= http.StatusBadRequest:
		reqCounter.Add(err500, 1)
	}

	reqCounter.Add(bytes, bc.responseBytes)
	reqCounter.Add(completed, 1)
	rps.Incr(1) // nolint:gomnd
	bps.Incr(bc.responseBytes)
	durAvg.Incr(dur.Nanoseconds())

	reqSecond.Set(rps.Rate())
	bytesSecond.Set(bps.Rate())
	avgDuration.Set(durAvg.Rate())
}
