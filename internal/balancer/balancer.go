package balancer

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Balancer struct {
	logger   *slog.Logger
	strategy Strategy
	proxies  map[string]*httputil.ReverseProxy
}

func NewBalancer(
	logger *slog.Logger,
	strategy Strategy,
	backends []string,
	urlCreateFunc func(string) *url.URL,
) *Balancer {
	proxies := make(map[string]*httputil.ReverseProxy, len(backends))
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

	b.proxies[backend].ServeHTTP(w, r)
}

type ErrorResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"status"`
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
