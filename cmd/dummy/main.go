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
	defer r.Body.Close()
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Warn("Read body",
			slog.String("error", err.Error()),
		)
	}

	h.logger.Info("Request",
		slog.String("Method", r.Method),
		slog.String("URI", r.RequestURI),
		slog.String("Body", string(bytes)),
	)
}

func main() {
	logger := slog.Default()

	server := http.Server{
		Addr:    "localhost:8081",
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
	}

	<-shutdownWaitChan
}
