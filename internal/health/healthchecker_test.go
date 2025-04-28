package health

import (
	"context"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/return_ok", func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()

		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/return_500", func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()

		w.WriteHeader(http.StatusServiceUnavailable)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("with response ok", func(t *testing.T) {
		t.Parallel()

		mockObserver := mock_observer.NewMockObserver(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())

		checker := NewChecker(
			slog.Default(),
			server.Client(),
			server.URL,
			func(s string) string { return server.URL + "/return_ok" },
			time.Millisecond*80,
			time.Millisecond*90,
			mockObserver,
		)

		mockObserver.EXPECT().UpdateBackendHealth(server.URL, true)

		go checker.Run(ctx)

		time.Sleep(time.Millisecond * 100)
		cancel()
	})

	t.Run("with response status 5xx", func(t *testing.T) {
		t.Parallel()

		mockObserver := mock_observer.NewMockObserver(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())

		checker := NewChecker(
			slog.Default(),
			server.Client(),
			server.URL,
			func(s string) string { return server.URL + "/return_500" },
			time.Millisecond*500,
			time.Millisecond*500,
			mockObserver,
		)

		mockObserver.EXPECT().UpdateBackendHealth(server.URL, false)

		go checker.Run(ctx)

		time.Sleep(time.Millisecond * 600)
		cancel()
	})
}
