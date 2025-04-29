package balancer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var errNoAvailableBackends = errors.New("no available backends")

// Balancer is reverse proxy that balance incoming requests between backends
// according to the given strategy.
type Balancer struct {
	logger   *slog.Logger
	strategy Strategy
	proxies  map[string]http.Handler
}

// NewBalancer creates Balancer.
func NewBalancer(
	logger *slog.Logger,
	strategy Strategy,
	backends []string,
	urlCreateFunc func(string) *url.URL,
) *Balancer {
	proxies := make(map[string]http.Handler, len(backends))
	for _, backend := range backends {
		rp := httputil.NewSingleHostReverseProxy(urlCreateFunc(backend))
		rp.ErrorHandler = createErrorHandler(logger.With(slog.String("backend", backend)))
		proxies[backend] = rp
	}

	return &Balancer{
		logger:   logger,
		strategy: strategy,
		proxies:  proxies,
	}
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := b.strategy.ChooseBackend()
	logger := b.logger.With(
		slog.String("method", r.Method),
		slog.String("url", r.RequestURI),
		slog.String("chosen_backend", backend),
	)

	if backend == "" {
		logger.Error("No available backends for request")
		writeErrorToClient(w, http.StatusServiceUnavailable, errNoAvailableBackends)
		return
	}

	proxy, ok := b.proxies[backend]
	if !ok {
		logger.Error("Unknown backend")
		writeErrorToClient(w, http.StatusInternalServerError, fmt.Errorf("strategy returned not existing backend: %s", backend))
		return
	}

	logger.Info("Serving request")

	proxy.ServeHTTP(w, r)
}

// ErrorResponse returned to client, when error occurred.
type ErrorResponse struct {
	// Msg includes occurred error.
	Msg string `json:"msg"`
	// Code is the returned http status code.
	Code int `json:"status"`
}

func createErrorHandler(logger *slog.Logger) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error("Error from backend",
			slog.String("error", err.Error()),
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
		)

		writeErrorToClient(w, http.StatusInternalServerError, fmt.Errorf("error from backend: %w", err))
	}
}

func writeErrorToClient(w http.ResponseWriter, statusCode int, err error) {
	dto := ErrorResponse{
		Msg:  err.Error(),
		Code: statusCode,
	}

	w.WriteHeader(dto.Code)

	encoder := json.NewEncoder(w)
	_ = encoder.Encode(dto)
}
