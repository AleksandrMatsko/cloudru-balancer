package health

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_observer "github.com/AleksandrMatsko/cloudru-balancer/internal/health/mocks"
	"go.uber.org/mock/gomock"
)

func TestChecker(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "return_ok" || r.URL.Path == "/return_ok" {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "healthy")
				return
			}

			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "not healthy")
		},
	))
	defer server.Close()

	t.Run("with response ok", func(t *testing.T) {
		mockObserver := mock_observer.NewMockObserver(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())

		checker := NewChecker(
			slog.Default(),
			server.Client(),
			server.URL,
			func(s string) string { return server.URL + "/return_ok" },
			time.Millisecond*100,
			time.Millisecond*100,
			mockObserver,
		)

		mockObserver.EXPECT().UpdateBackendHealth(server.URL, true).Times(1)

		go checker.Run(ctx)

		time.Sleep(time.Millisecond * 120)

		cancel()
	})

	t.Run("with response status 5xx", func(t *testing.T) {
		mockObserver := mock_observer.NewMockObserver(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())

		checker := NewChecker(
			slog.Default(),
			server.Client(),
			server.URL,
			func(s string) string { return server.URL + "/return_5xx" },
			time.Millisecond*100,
			time.Millisecond*100,
			mockObserver,
		)

		mockObserver.EXPECT().UpdateBackendHealth(server.URL, false).Times(1)

		go checker.Run(ctx)

		time.Sleep(time.Millisecond * 150)

		cancel()
	})
}
