package balancer

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

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
	b.logger.Info("Request",
		slog.String("method", r.Method),
		slog.String("url", r.RequestURI),
		slog.String("chosen_backend", backend),
	)

	if backend == "" {
		dto := ErrorResponse{
			Msg:  "no available backends",
			Code: http.StatusServiceUnavailable,
		}

		w.WriteHeader(dto.Code)

		encoder := json.NewEncoder(w)
		_ = encoder.Encode(dto)

		return
	}

	proxy, ok := b.proxies[backend]
	if !ok {
		dto := ErrorResponse{
			Msg:  fmt.Sprintf("strategy returned not existing backend: %s", backend),
			Code: http.StatusInternalServerError,
		}

		w.WriteHeader(dto.Code)

		encoder := json.NewEncoder(w)
		_ = encoder.Encode(dto)

		return
	}

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

		dto := ErrorResponse{
			Msg:  err.Error(),
			Code: http.StatusInternalServerError,
		}

		w.WriteHeader(dto.Code)

		encoder := json.NewEncoder(w)
		_ = encoder.Encode(dto)
	}
}
