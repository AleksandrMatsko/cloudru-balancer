package balancer

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_balancer "github.com/AleksandrMatsko/cloudru-balancer/internal/balancer/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBalancer_ServeHTTP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("when strategy returns empty backend", func(t *testing.T) {
		t.Parallel()

		mockStrategy := mock_balancer.NewMockStrategy(mockCtrl)

		b := Balancer{
			logger:   slog.Default(),
			strategy: mockStrategy,
		}

		mockStrategy.EXPECT().ChooseBackend().Return("").Times(1)
		expectedDTO := ErrorResponse{
			Msg:  "no available backends",
			Code: http.StatusServiceUnavailable,
		}

		expectedBytes, err := json.Marshal(expectedDTO)
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://test.url", nil)

		b.ServeHTTP(recorder, req)

		assert.Equal(t, expectedDTO.Code, recorder.Code)

		bytes, err := io.ReadAll(recorder.Body)
		assert.Nil(t, err)
		assert.Equal(t, string(expectedBytes)+"\n", string(bytes))
	})

	t.Run("when strategy returns non existed backend", func(t *testing.T) {
		t.Parallel()

		mockStrategy := mock_balancer.NewMockStrategy(mockCtrl)

		b := Balancer{
			logger:   slog.Default(),
			strategy: mockStrategy,
		}

		mockStrategy.EXPECT().ChooseBackend().Return("hello").Times(1)
		expectedDTO := ErrorResponse{
			Msg:  "strategy returned not existing backend: hello",
			Code: http.StatusInternalServerError,
		}

		expectedBytes, err := json.Marshal(expectedDTO)
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://test.url", nil)

		b.ServeHTTP(recorder, req)

		assert.Equal(t, expectedDTO.Code, recorder.Code)

		bytes, err := io.ReadAll(recorder.Body)
		assert.Nil(t, err)
		assert.Equal(t, string(expectedBytes)+"\n", string(bytes))
	})

	t.Run("when strategyr returns valid backend", func(t *testing.T) {
		t.Parallel()

		mockStrategy := mock_balancer.NewMockStrategy(mockCtrl)
		mockProxy := mock_balancer.NewMockHandler(mockCtrl)

		backendHost := "my.backend.com"

		b := Balancer{
			logger:   slog.Default(),
			strategy: mockStrategy,
			proxies: map[string]http.Handler{
				backendHost: mockProxy,
			},
		}

		mockStrategy.EXPECT().ChooseBackend().Return("hello").Times(1)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://test.url", nil)

		mockProxy.EXPECT().ServeHTTP(recorder, req).Times(1)

		b.ServeHTTP(recorder, req)
	})
}
