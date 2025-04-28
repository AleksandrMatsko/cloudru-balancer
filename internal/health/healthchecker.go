// health contains entity that checks backend availability.
package health

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Checker checks if the backend with specified url is healthy by sending GET requests.
// Backend is healhy when:
//   - no connection problems occurred;
//   - response on check health request does not have code in range [500, 599] inclusive.
type Checker struct {
	logger         *slog.Logger
	client         *http.Client
	backend        string
	urlCreateFunc  func(string) string
	checkTimeout   time.Duration
	requestTimeout time.Duration
	observer       Observer
	wasHealfy      bool
}

// NewChecker creates new checker.
func NewChecker(
	logger *slog.Logger,
	client *http.Client,
	backend string,
	urlCreateFunc func(string) string,
	checkTimeout, requestTimeout time.Duration,
	observer Observer,
) *Checker {
	return &Checker{
		logger:         logger.With(slog.String("backend_healthcheck_url", urlCreateFunc(backend))),
		client:         client,
		backend:        backend,
		urlCreateFunc:  urlCreateFunc,
		checkTimeout:   checkTimeout,
		requestTimeout: requestTimeout,
		observer:       observer,
		wasHealfy:      false,
	}
}

// Run check loop. Should be started in separate goroutine.
func (checker *Checker) Run(ctx context.Context) {
	ticker := time.NewTicker(checker.checkTimeout)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			available := checker.check(ctx)
			checker.observer.UpdateBackendHealth(checker.backend, available)
		}
	}
}

func (checker *Checker) check(ctx context.Context) bool {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, checker.requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, checker.urlCreateFunc(checker.backend), nil)
	if err != nil {
		checker.logger.Warn("Create healthcheck request",
			slog.String("error", err.Error()),
		)
		return false
	}

	rsp, err := checker.client.Do(req)
	if err != nil {
		if checker.wasHealfy {
			checker.wasHealfy = false
			checker.logger.Warn("Backend unavailable",
				slog.String("error", err.Error()),
			)
		}

		return false
	}

	defer rsp.Body.Close()
	_, _ = io.ReadAll(rsp.Body)

	if is5xxCode(rsp.StatusCode) {
		if checker.wasHealfy {
			checker.wasHealfy = false
			checker.logger.Warn("Backend returns 5xx status code",
				slog.String("status", rsp.Status),
			)
		}

		return false
	}

	if !checker.wasHealfy {
		checker.wasHealfy = true
		checker.logger.Info("Backend is available again")
	}

	return true
}

func is5xxCode(code int) bool {
	return code >= 500 && code < 600
}
