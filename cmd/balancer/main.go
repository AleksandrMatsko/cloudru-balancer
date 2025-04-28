package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/AleksandrMatsko/cloudru-balancer/internal/balancer"
	"github.com/AleksandrMatsko/cloudru-balancer/internal/config"
	"github.com/AleksandrMatsko/cloudru-balancer/internal/health"
	"github.com/AleksandrMatsko/cloudru-balancer/internal/strategies"
)

var (
	configFileNameFlag = flag.String("config", "/etc/cloudru_balancer/balancer.yml", "Path to configuration file")
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

	strategy, err := createStrategy(appConfig)
	if err != nil {
		logger.Error("Select balancing strategy",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	runHealthCheckers(ctx, logger, appConfig, strategy)

	balancer := balancer.NewBalancer(
		logger,
		strategy,
		appConfig.Backends,
		createURL,
	)

	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", appConfig.Port),
		Handler: balancer,
	}

	shutdownWaitChan := make(chan os.Signal)

	go func() {
		sigWaitChan := make(chan os.Signal, 1)
		signal.Notify(sigWaitChan, os.Interrupt)

		<-sigWaitChan

		cancel()
		if err := server.Shutdown(context.TODO()); err != nil {
			logger.Warn("Shutdown",
				slog.String("error", err.Error()))
		}
		close(shutdownWaitChan)
	}()

	logger.Info("Listen",
		slog.String("address", server.Addr),
	)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Warn("ListenAndServe",
			slog.String("error", err.Error()))
	}

	<-shutdownWaitChan
}

type observingStrategy interface {
	health.Observer
	balancer.Strategy
}

func createStrategy(conf config.Balancer) (observingStrategy, error) {
	switch conf.Strategy {
	case "RoundRobin":
		return strategies.NewRoundRobin(
			conf.Backends,
		), nil
	default:
		return nil, fmt.Errorf("unknown strategy: %s", conf.Strategy)
	}
}

func runHealthCheckers(ctx context.Context, logger *slog.Logger, conf config.Balancer, observer health.Observer) {
	client := &http.Client{}

	for _, backend := range conf.Backends {
		checker := health.NewChecker(
			logger,
			client,
			backend,
			createURLString,
			time.Duration(conf.Heathcheck.CheckTimeoutSeconds)*time.Second,
			time.Duration(conf.Heathcheck.RequestTimeoutSeconds)*time.Second,
			observer,
		)

		go checker.Run(ctx)
	}
}

func createURL(backend string) *url.URL {
	url, _ := url.Parse(createURLString(backend))
	return url
}

func createURLString(backend string) string {
	return fmt.Sprintf("http://%s", backend)
}
