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

	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
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

func loggingMiddleware() func(http.Handler) http.Handler {
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
				slog.Error("HTTP Request Error", request, response)
			case lrw.statusCode >= 400:
				slog.Warn("HTTP Request Warning", request, response)
			default:
				slog.Info("HTTP Request", request, response)
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

func addRoutes(mux *http.ServeMux, c *cache.Cache, group *singleflight.Group) {
	mux.HandleFunc("GET /", AppConfigHandler)
	mux.HandleFunc("GET /logo", LogoHandler)

	mux.HandleFunc("POST /oauth2/v1/authorize", AuthorizeHandler)
	mux.HandleFunc("POST /oauth2/v1/access_token", TokenHandler)
	mux.HandleFunc("POST /validate", ValidateHandler)

	mux.HandleFunc("POST /api/v1/synchronizer/config", SyncConfigHandler)
	mux.HandleFunc("POST /api/v1/synchronizer/schema", SchemaHandler)
	mux.HandleFunc("POST /api/v1/synchronizer/data", DataHandler(c, group))
	mux.HandleFunc("POST /api/v1/synchronizer/filter/validate", ValidateFiltersHandler)

	mux.Handle("POST /api/v1/automations/sync_action/{type}", ActionAuth(http.HandlerFunc(ActionHandler)))

	mux.HandleFunc("POST /api/v1/synchronizer/webhooks", RegisterHandler)
	mux.HandleFunc("POST /api/v1/synchronizer/webhooks/verify", VerifyHandler)
	mux.HandleFunc("POST /api/v1/synchronizer/webhooks/transform", TransformHandler)
	mux.HandleFunc("DELETE /api/v1/synchronizer/webhooks", DeleteHandler)
}

func NewServer(c *cache.Cache, group *singleflight.Group) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, c, group)
	var handler http.Handler = mux

	handler = loggingMiddleware()(handler)
	handler = gzipMiddleware(handler)
	return handler
}
