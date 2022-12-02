IMG ?= xakdwch5/etcd-carry:latest
PLATFORMS ?= linux/amd64,linux/arm64

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: build
build: fmt vet
	CGO_ENABLED=1 go build -a -o bin/etcd-carry main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: golangci-lint
	$(GOLANGCI_LINT) run

test: fmt vet
	go test ./pkg/... -coverprofile cover.out

docker-build:
	@echo "build docker image ${IMG}"
	@docker build --pull --no-cache . -t ${IMG}

docker-push:
	@echo "push docker image ${IMG}"
	@docker push ${IMG}

docker-multiarch:
	docker buildx build -f ./Dockerfile_multiarch --pull --no-cache --platform=$(PLATFORMS) --push . -t $(IMG)

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
golangci-lint:
	 $(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1)

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
