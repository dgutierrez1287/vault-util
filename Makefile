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
	@echo "install - installs the executable"
	@echo "uninstall - uninstalls the execuable"

build:
	go build .

install:
	sudo cp -f ./vault-util $(INSTALL_DIR) && sudo chmod 777 $(INSTALL_DIR)/vault-util 

uninstall:
	sudo rm -f $(INSTALL_DIR)/vault-util

tidy:
	go mod tidy
