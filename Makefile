GO_PATH := $(shell go env GOPATH)
GOLANGCI_LINT_VERSION := ""

.PHONY: build-balancer
build-balancer:
	CGO_ENABLED=0 GOOS=linux go build github.com/AleksandrMatsko/cloudru-balancer/cmd/balancer

.PHONY: build-dummy
build-dummy:
	CGO_ENABLED=0 GOOS=linux go build github.com/AleksandrMatsko/cloudru-balancer/cmd/dummy

.PHONY: test
test:
	go test -v -bench=. -race ./...

.PHONY: install-lint
install-lint:
	# The recommended way to install golangci-lint into CI/CD
	wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GO_PATH}/bin ${GOLANGCI_LINT_VERSION}

.PHONY: lint
lint:
	golangci-lint run

.PHONY: mock
mock:
	./scripts/generate_mocks.sh
