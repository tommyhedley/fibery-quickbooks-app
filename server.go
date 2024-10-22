package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/tommyhedley/fibery/fibery-tsheets-integration/actions"
	"github.com/tommyhedley/fibery/fibery-tsheets-integration/oauth2"
	"github.com/tommyhedley/fibery/fibery-tsheets-integration/synchronizer"
	"github.com/tommyhedley/fibery/fibery-tsheets-integration/webhooks"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	if lrw.statusCode >= 400 {
		lrw.body.Write(b)
	}
	return n, err
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func loggingMiddleware(httpLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			request := slog.Group(
				"Request",
				"Method", r.Method,
				"Path", r.URL.Path,
				"CorrelationID", r.Header.Get("X-Correlationid"),
				"ClientIP", r.RemoteAddr,
				"Duration", duration.String(),
			)

			response := slog.Group("Response", "StatusCode", lrw.statusCode)

			if lrw.statusCode >= 400 && lrw.body.Len() > 0 {
				var errorResponse struct {
					Error    string `json:"error"`
					TryLater bool   `json:"tryLater"`
				}
				err := json.Unmarshal(lrw.body.Bytes(), &errorResponse)
				if err != nil {
					response = slog.Group("Response", "StatusCode", lrw.statusCode, "ErrorMsg", lrw.body.String())
				} else {
					response = slog.Group("Response", "StatusCode", lrw.statusCode, "ErrorMsg", errorResponse.Error)
					if errorResponse.TryLater {
						response = slog.Group("Response", "StatusCode", lrw.statusCode, "ErrorMsg", errorResponse.Error, "TryLater", true)
					}
				}
			}

			switch {
			case lrw.statusCode >= 500:
				httpLogger.Error("HTTP Request Error", request, response)
			case lrw.statusCode >= 400:
				httpLogger.Warn("HTTP Request Warning", request, response)
			default:
				httpLogger.Info("HTTP Request", request, response)
			}
		})
	}
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzrw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}

		next.ServeHTTP(gzrw, r)
	})
}

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", ConfigHandler)
	mux.HandleFunc("GET /logo", LogoHandler)

	mux.HandleFunc("POST /oauth2/v1/authorize", oauth2.AuthorizeHandler)
	mux.HandleFunc("POST /oauth2/v1/access_token", oauth2.TokenHandler)
	mux.HandleFunc("POST /validate", oauth2.ValidateHandler)

	mux.HandleFunc("POST /api/v1/synchronizer/config", synchronizer.ConfigHandler)
	mux.HandleFunc("POST /api/v1/synchronizer/schema", synchronizer.SchemaHandler)
	mux.HandleFunc("POST /api/v1/synchronizer/data", synchronizer.DataHandler)
	mux.HandleFunc("POST /api/v1/synchronizer/filter/validate", synchronizer.ValidateFiltersHandler)

	mux.Handle("POST /api/v1/automations/sync_action/{type}", actions.SyncActionAuth(http.HandlerFunc(actions.SyncActionHandler)))

	mux.HandleFunc("POST /api/v1/synchronizer/webhooks", webhooks.RegisterHandler)
	mux.HandleFunc("POST /api/v1/synchronizer/webhooks/verify", webhooks.VerifyHandler)
	mux.HandleFunc("POST /api/v1/synchronizer/webhooks/transform", webhooks.TransformHandler)
	mux.HandleFunc("DELETE /api/v1/synchronizer/webhooks", webhooks.DeleteHandler)
}

func NewServer(httpLogger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)
	var handler http.Handler = mux

	handler = loggingMiddleware(httpLogger)(handler)
	handler = gzipMiddleware(handler)
	return handler
}
