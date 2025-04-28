package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/AleksandrMatsko/cloudru-balancer/internal/config"
)

var (
	configFileNameFlag = flag.String("config", "./cloudru_balancer.yml", "Path to configuration file")
	printConfigFlag    = flag.Bool("print-config", false, "Print current config to stdout")
)

func main() {
	flag.Parse()
	logger := slog.Default()

	appConfig := config.DefaultForBalancer()
	err := config.Read(*configFileNameFlag, &appConfig)
	if err != nil {
		logger.Error("Read config",
			slog.String("config_file_name", *configFileNameFlag),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	if *printConfigFlag {
		config.Print(appConfig)
	}

	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", appConfig.Port),
		Handler: nil,
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
