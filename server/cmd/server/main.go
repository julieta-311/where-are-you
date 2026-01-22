package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/julieta-311/where-are-you/server/internal/handlers"
	"github.com/julieta-311/where-are-you/server/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize OpenTelemetry Tracer
	// In production, you would likely send this to a collector (e.g., Jaeger/OTLP).
	// Here we write to stdout for demonstration.
	shutdown, err := telemetry.InitTracer(os.Stdout)
	if err != nil {
		logger.Error("failed to init tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			logger.Error("failed to shutdown tracer", "error", err)
		}
	}()

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "whereareyou.db"
	}

	store, err := handlers.NewSQLiteStore(dbPath)
	if err != nil {
		logger.Error("failed to init sqlite store", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	h := handlers.NewHandler(store, logger)

	r := mux.NewRouter()
	r.Use(otelhttp.NewMiddleware("server"))
	r.Use(loggingMiddleware(logger))

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/requests", h.CreateRequest).Methods(http.MethodPost)
	api.HandleFunc("/requests", h.GetPendingRequests).Methods(http.MethodGet)
	api.HandleFunc("/requests/{requestId}/respond", h.RespondToRequest).Methods(http.MethodPost)
	api.HandleFunc("/location", h.UpdateLocation).Methods(http.MethodPost)
	api.HandleFunc("/location/{userId}", h.GetLocation).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen and serve failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("forced shutdown", "error", err)
	}
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)

			span := trace.SpanFromContext(r.Context())
			traceID := span.SpanContext().TraceID().String()
			spanID := span.SpanContext().SpanID().String()

			logger.Info("request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
				"trace_id", traceID,
				"span_id", spanID,
			)
		})
	}
}
