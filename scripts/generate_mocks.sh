#!/bin/bash

go install go.uber.org/mock/mockgen@v0.5.1

rm -r ./internal/health/mocks/*

mockgen -destination=internal/health/mocks/observer.go -package=mock_observer github.com/AleksandrMatsko/cloudru-balancer/internal/health Observer
