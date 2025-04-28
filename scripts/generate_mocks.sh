#!/bin/bash

go install go.uber.org/mock/mockgen@v0.5.1

rm -r ./internal/health/mocks/*

mockgen -destination=internal/health/mocks/observer.go -package=mock_observer github.com/AleksandrMatsko/cloudru-balancer/internal/health Observer

rm -r ./internal/balancer/mocks/*

mockgen -destination=internal/balancer/mocks/strategy.go -package=mock_balancer github.com/AleksandrMatsko/cloudru-balancer/internal/balancer Strategy
mockgen -destination=internal/balancer/mocks/http_handler.go -package=mock_balancer net/http Handler