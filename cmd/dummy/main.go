package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

type handler struct {
	logger *slog.Logger
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(
		slog.String("method", r.Method),
		slog.String("uri", r.RequestURI),
	)

	defer r.Body.Close()
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Warn("Read body",
			slog.String("error", err.Error()),
		)
	}

	logger.Info("Request",
		slog.String("body", string(bytes)),
	)
}

func main() {
	logger := slog.Default()

	server := http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: &handler{logger: logger},
	}

	shutdownWaitChan := make(chan os.Signal)

	go func() {
		sigWaitChan := make(chan os.Signal, 1)
		signal.Notify(sigWaitChan, os.Interrupt)

		<-sigWaitChan

		if err := server.Shutdown(context.TODO()); err != nil {
			logger.Warn("Shutdown",
				slog.String("error", err.Error()))
		}
		close(shutdownWaitChan)
	}()

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Warn("ListenAndServe",
			slog.String("error", err.Error()))
		return
	}

	<-shutdownWaitChan
}
