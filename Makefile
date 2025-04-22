MAKEFILE_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

# Detect OS
OS := $(shell uname)

INSTALL_DIR ?= "/usr/local/bin"

help:
	@echo "Vault-Util Makefile"
	@echo "==============================="
	@echo 
	@echo "help - shows all the help information"
	@echo "build - builds the vault-loader executable"
	@echo "clean - cleans the directory of un-needed files"
	@echo "install - installs the executable"
	@echo "uninstall - uninstalls the execuable"
	@echo "tidy - runs go mod tidy to update modules"
	@echo "go-fmt - runs go fmt on all the files"
	@echo "test - runs all tests for the project"
	@echo "test-package <pkg_name> - runs tests for only a single package"
	@echo "coverage - runs tests and outputs the test coverage report"

build:
	go build .

clean:
	rm -f /.vault-util
	rm -f coverage.out coverage.html

install:
	sudo cp -f ./vault-util $(INSTALL_DIR) && sudo chmod 777 $(INSTALL_DIR)/vault-util 

uninstall:
	sudo rm -f $(INSTALL_DIR)/vault-util

tidy:
	go mod tidy

go-fmt:
	go fmt -v ./...

test:
	go test -v ./...

test-package:
	@if [ -z "${pkg}" ]; then \
		echo "Usage: make test-package pkg=<package-name>"; \
		exit 1; \
	fi

	go test -v ./${pkg}

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
