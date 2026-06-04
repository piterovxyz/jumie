.PHONY: build clean install uninstall

OS := $(shell uname -s)

build:
	go build -o dist/jum ./cmd/jumie
	go build -o dist/jumied ./cmd/jumied

clean:
	rm -rf dist/

install: build
	mkdir -p $(HOME)/.local/bin
	cp dist/jum $(HOME)/.local/bin/
	cp dist/jumied $(HOME)/.local/bin/
	chmod +x $(HOME)/.local/bin/jum $(HOME)/.local/bin/jumied
ifeq ($(OS),Darwin)
	mkdir -p $(HOME)/Library/LaunchAgents
	sed "s|HOME_DIR|$(HOME)|g" org.jumie.jumied.plist > $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist
	launchctl unload $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist 2>/dev/null || true
	launchctl load $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist
	@echo "daemon registered in launchd. socket will be available at ~/.local/share/jumie/jumie.sock"
else ifeq ($(OS),Linux)
	mkdir -p $(HOME)/.config/systemd/user
	sed "s|HOME_DIR|$(HOME)|g" jumied.service > $(HOME)/.config/systemd/user/jumied.service
	systemctl --user daemon-reload
	systemctl --user enable --now jumied
	@echo "daemon registered in systemd. socket will be available at ~/.local/share/jumie/jumie.sock"
endif

uninstall:
ifeq ($(OS),Darwin)
	launchctl unload $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist 2>/dev/null || true
	rm -f $(HOME)/Library/LaunchAgents/org.jumie.jumied.plist
else ifeq ($(OS),Linux)
	systemctl --user disable --now jumied 2>/dev/null || true
	rm -f $(HOME)/.config/systemd/user/jumied.service
	systemctl --user daemon-reload
endif
	rm -f $(HOME)/.local/bin/jum $(HOME)/.local/bin/jumied
	@echo "jumie uninstalled successfully."
