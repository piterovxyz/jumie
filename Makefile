SHELL := /bin/bash

color_reset   := \033[0m
color_error   := \033[31m
color_success := \033[32m
color_info    := \033[36m
color_warn    := \033[33m

OS := $(shell uname -s)

.PHONY: build clean install uninstall help check-go check-systemd release

.DEFAULT_GOAL := help

help: ## show help details
	@printf "$(color_info)available commands:$(color_reset)\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(color_success)%-15s$(color_reset) %s\n", $$1, $$2}'

check-go:
	@which go > /dev/null 2>&1 || (printf "$(color_error)error: go is not installed. please install go before building.$(color_reset)\n" && exit 1)

check-systemd:
ifeq ($(OS),Linux)
	@systemctl --user show-environment >/dev/null 2>&1 || \
		(printf "$(color_warn)warning: systemctl --user is not available. dbus might not be running or XDG_RUNTIME_DIR is unset.$(color_reset)\n" && \
		 printf "$(color_warn)try: export XDG_RUNTIME_DIR=/run/user/\$$(id -u)$(color_reset)\n")
endif

build: check-go ## build jumie and jumied binaries into dist
	@printf "$(color_info)building binaries...$(color_reset)\n"
	@mkdir -p dist
	@go build -o dist/jum ./cmd/jumie
	@go build -o dist/jumied ./cmd/jumied
	@printf "$(color_success)build completed successfully! binaries are in dist/$(color_reset)\n"

release: check-go ## build multi-platform release archives into dist/releases
	@printf "$(color_info)building releases...$(color_reset)\n"
	@rm -rf dist/releases
	@mkdir -p dist/releases
	@GOOS=linux GOARCH=amd64 go build -o dist/releases/linux-amd64/jum ./cmd/jumie
	@GOOS=linux GOARCH=amd64 go build -o dist/releases/linux-amd64/jumied ./cmd/jumied
	@cp jumied.service dist/releases/linux-amd64/
	@COPYFILE_DISABLE=1 tar -czf dist/releases/jumie-linux-amd64.tar.gz -C dist/releases/linux-amd64 jum jumied jumied.service
	@GOOS=linux GOARCH=arm64 go build -o dist/releases/linux-arm64/jum ./cmd/jumie
	@GOOS=linux GOARCH=arm64 go build -o dist/releases/linux-arm64/jumied ./cmd/jumied
	@cp jumied.service dist/releases/linux-arm64/
	@COPYFILE_DISABLE=1 tar -czf dist/releases/jumie-linux-arm64.tar.gz -C dist/releases/linux-arm64 jum jumied jumied.service
	@GOOS=darwin GOARCH=amd64 go build -o dist/releases/darwin-amd64/jum ./cmd/jumie
	@GOOS=darwin GOARCH=amd64 go build -o dist/releases/darwin-amd64/jumied ./cmd/jumied
	@cp org.jumie.jumied.plist dist/releases/darwin-amd64/
	@COPYFILE_DISABLE=1 tar -czf dist/releases/jumie-darwin-amd64.tar.gz -C dist/releases/darwin-amd64 jum jumied org.jumie.jumied.plist
	@GOOS=darwin GOARCH=arm64 go build -o dist/releases/darwin-arm64/jum ./cmd/jumie
	@GOOS=darwin GOARCH=arm64 go build -o dist/releases/darwin-arm64/jumied ./cmd/jumied
	@cp org.jumie.jumied.plist dist/releases/darwin-arm64/
	@COPYFILE_DISABLE=1 tar -czf dist/releases/jumie-darwin-arm64.tar.gz -C dist/releases/darwin-arm64 jum jumied org.jumie.jumied.plist
	@rm -rf dist/releases/linux-amd64 dist/releases/linux-arm64 dist/releases/darwin-amd64 dist/releases/darwin-arm64
	@printf "$(color_success)releases built successfully in dist/releases/$(color_reset)\n"

clean: ## remove built binaries and temporary artifacts
	@printf "$(color_info)cleaning dist/...$(color_reset)\n"
	@rm -rf dist/
	@printf "$(color_success)cleaned!$(color_reset)\n"

install: build check-systemd ## install binaries and register the daemon
	@printf "$(color_info)installing binaries to ~/.local/bin/...$(color_reset)\n"
	@mkdir -p $(HOME)/.local/bin
	@cp dist/jum $(HOME)/.local/bin/
	@cp dist/jumied $(HOME)/.local/bin/
	@chmod +x $(HOME)/.local/bin/jum $(HOME)/.local/bin/jumied
ifeq ($(OS),Darwin)
	@printf "$(color_info)configuring launchd for macos...$(color_reset)\n"
	@mkdir -p $(HOME)/Library/LaunchAgents
	@sed "s|HOME_DIR|$(HOME)|g" org.jumie.jumied.plist > $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist
	@launchctl unload $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist 2>/dev/null || true
	@launchctl load $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist
	@printf "$(color_success)installation completed! socket will be available at ~/.local/share/jumie/jumie.sock$(color_reset)\n"
else ifeq ($(OS),Linux)
	@printf "$(color_info)configuring systemd for linux...$(color_reset)\n"
	@mkdir -p $(HOME)/.config/systemd/user
	@sed "s|HOME_DIR|$(HOME)|g" jumied.service > $(HOME)/.config/systemd/user/jumied.service
	@if systemctl --user show-environment >/dev/null 2>&1; then \
		systemctl --user daemon-reload; \
		systemctl --user enable --now jumied; \
		printf "$(color_success)installation completed! jumied daemon is running.$(color_reset)\n"; \
	else \
		printf "$(color_warn)warning: service file copied, but not started because systemctl --user is unavailable.$(color_reset)\n"; \
	fi
	@printf "$(color_success)socket will be available at ~/.local/share/jumie/jumie.sock$(color_reset)\n"
endif

uninstall: ## uninstall binaries and unregister the daemon
	@printf "$(color_info)removing installed components...$(color_reset)\n"
ifeq ($(OS),Darwin)
	@launchctl unload $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist 2>/dev/null || true
	@rm -f $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist
else ifeq ($(OS),Linux)
	@if systemctl --user show-environment >/dev/null 2>&1; then \
		systemctl --user disable --now jumied 2>/dev/null || true; \
	fi
	@rm -f $(HOME)/.config/systemd/user/jumied.service
	@if systemctl --user show-environment >/dev/null 2>&1; then \
		systemctl --user daemon-reload; \
	fi
endif
	@rm -f $(HOME)/.local/bin/jum $(HOME)/.local/bin/jumied
	@printf "$(color_success)uninstallation completed!$(color_reset)\n"
