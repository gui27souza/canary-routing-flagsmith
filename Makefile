-include .env
export

export CGO_ENABLED := 0

# OS Detection
ifeq ($(OS),Windows_NT)
	BINARY := bin\api.exe
	CLEAN_CMD := if exist bin rmdir /s /q bin && if exist coverage.out del /f /q coverage.out
else
	BINARY := ./bin/api
	CLEAN_CMD := rm -rf bin/ coverage.out
endif

DOCKER_IMAGE := go-flagsmith:latest

.DEFAULT_GOAL := build

.PHONY: fmt vet test coverage check-env dev build start clean docker-build docker-run load-test

fmt:
	@echo Formatting...
	@go fmt ./...

vet: fmt
	@echo Vetting...
	@go vet ./...

test:
	@echo Testing...
	@go test -v ./...

ALL_PKGS := $(shell go list ./internal/...)
FILTERED_PKGS := $(filter-out %/testutil, $(ALL_PKGS))
COMMA := ,
EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
COVER_PKGS := $(subst $(SPACE),$(COMMA),$(FILTERED_PKGS))

coverage:
	@echo Testing and checking coverage...
	@go test -coverprofile=coverage.out -coverpkg=$(COVER_PKGS) ./internal/...
	@go tool cover -html=coverage.out

check-env:
	@$(if $(strip $(FLAGSMITH_API_KEY)),echo ✅ Variable FLAGSMITH_API_KEY is filled.,$(error ⚠️ Variable FLAGSMITH_API_KEY is empty!))

dev: vet check-env
	@echo Running code...
	@go run ./cmd/api

build: vet
	@echo Building...
	@go build -ldflags="-s -w" -o $(BINARY) ./cmd/api

start: build
	@echo Running binary artifact...
	@$(BINARY)

clean:
	@echo Cleaning artifacts...
	@$(CLEAN_CMD)

docker-build:
	@echo Building distroless docker image...
	@docker build -t $(DOCKER_IMAGE) .

docker-run: check-env
	@echo Running container on :8080...
	@docker run --rm -p 8080:8080 -e FLAGSMITH_API_KEY=$(FLAGSMITH_API_KEY) $(DOCKER_IMAGE)

load-test:
	@echo Running k6 load test...
	@k6 run load-test.js
