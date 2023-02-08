.DEFAULT_GOAL := help

SHELL := bash
VERSION := 0.3.0

# Load .env file if it exists.
ifneq (,$(wildcard ./.env))
  include .env
  export
endif


.PHONY: help
help: ## Show help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[/0-9a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'


.PHONY: format
format: ## Format source code
	@go fmt ./...


.PHONY: run
run: ## Run the app in dev mode
	@go run .


.PHONY: test
test: ## Run tests
	@go test -race -timeout 30m -cover ./...


.PHONY: test/verbos
test/verbose: ## Run tests with verbose outputting.
	@go test -race -timeout 30m -cover -v ./...


.PHONY: build/app
build/app: guard-ARCH
	$(eval APP_PATH := dist/apps/$(ARCH)/CloudSQLProxyMenuBar.app)
	$(eval INFO_PLIST_FILE := $(APP_PATH)/Contents/Info.plist)

	@echo "Building OSX app: $(APP_PATH) ..."
	@rm -rf $(APP_PATH)
	@mkdir -p $(APP_PATH)/Contents/{MacOS,Resources}
	@cp -p AppIcon.icns $(APP_PATH)/Contents/Resources/AppIcon.icns

	@echo '<?xml version="1.0" encoding="UTF-8"?>' > $(INFO_PLIST_FILE)
	@echo '<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' >> $(INFO_PLIST_FILE)
	@echo '<plist version="1.0">' >> $(INFO_PLIST_FILE)
	@echo '  <dict>' >> $(INFO_PLIST_FILE)
	@echo '    <key>CFBundleDevelopmentRegion</key>' >> $(INFO_PLIST_FILE)
	@echo '    <string>English</string>' >> $(INFO_PLIST_FILE)

	@echo '    <key>CFBundleExecutable</key>' >> $(INFO_PLIST_FILE)
	@echo '    <string>cloudsqlproxymenubar</string>' >> $(INFO_PLIST_FILE)

	@echo '    <key>CFBundleIconFile</key>' >> $(INFO_PLIST_FILE)
	@echo '    <string>AppIcon</string>' >> $(INFO_PLIST_FILE)

	@echo '    <key>CFBundleIconName</key>' >> $(INFO_PLIST_FILE)
	@echo '    <string>AppIcon</string>' >> $(INFO_PLIST_FILE)

	@echo '    <key>CFBundleIdentifier</key>' >> $(INFO_PLIST_FILE)
	@echo '    <string>dev.kohkimakimoto.CloudSQLProxyMenuBar</string>' >> $(INFO_PLIST_FILE)

	@echo '    <key>CFBundleName</key>' >> $(INFO_PLIST_FILE)
	@echo '    <string>CloudSQLProxyMenuBar</string>' >> $(INFO_PLIST_FILE)

	@echo '    <key>CFBundleSupportedPlatforms</key>' >> $(INFO_PLIST_FILE)
	@echo '    <array>' >> $(INFO_PLIST_FILE)
	@echo '      <string>MacOSX</string>' >> $(INFO_PLIST_FILE)
	@echo '    </array>' >> $(INFO_PLIST_FILE)

	@echo '    <key>CFBundlePackageType</key>' >> $(INFO_PLIST_FILE)
	@echo '    <string>APPL</string>' >> $(INFO_PLIST_FILE)

	@echo '    <key>CFBundleVersion</key>' >> $(INFO_PLIST_FILE)
	@echo '    <key>$(VERSION)</key>' >> $(INFO_PLIST_FILE)

	@echo '    <key>NSHighResolutionCapable</key>' >> $(INFO_PLIST_FILE)
	@echo '    <string>True</string>' >> $(INFO_PLIST_FILE)

	@echo '    <key>LSUIElement</key>' >> $(INFO_PLIST_FILE)
	@echo '    <string>1</string>' >> $(INFO_PLIST_FILE)
	@echo '  </dict>' >> $(INFO_PLIST_FILE)
	@echo '</plist>' >> $(INFO_PLIST_FILE)

	@echo "Building $(ARCH) executable ..."
	@GOOS=darwin GOARCH=$(ARCH) CGO_ENABLED=1 go build -ldflags="-s -w" -o="$(APP_PATH)/Contents/MacOS/cloudsqlproxymenubar" .


.PHONY: build/app/amd64
build/app/amd64: ## Build amd64 OSX app
	@$(MAKE) build/app ARCH=amd64


.PHONY: build/app/arm64
build/app/arm64: ## Build arm64 OSX app
	@$(MAKE) build/app ARCH=arm64


.PHONY: build/app/all
build/app/all: ## Build all OSX apps
	@$(MAKE) build/app/amd64
	@$(MAKE) build/app/arm64


.PHONY: build/dmg
build/dmg: guard-ARCH
	@$(MAKE) build/app ARCH=$(ARCH)

	$(eval DMG_PATH := dist/dmg/$(ARCH))
	$(eval APPDMG_FILE := $(DMG_PATH)/appdmg.json)

	@rm -rf $(DMG_PATH)
	@mkdir -p $(DMG_PATH)

	@echo '{' > $(APPDMG_FILE)
	@echo '  "title": "CloudSQLProxyMenuBarInstall",' >> $(APPDMG_FILE)
	@echo '  "icon": "../../../AppIcon.icns",' >> $(APPDMG_FILE)
	@echo '  "contents": [' >> $(APPDMG_FILE)
	@echo '    { "x": 448, "y": 344, "type": "link", "path": "/Applications" },' >> $(APPDMG_FILE)
	@echo '    { "x": 192, "y": 344, "type": "file", "path": "../../apps/$(ARCH)/CloudSQLProxyMenuBar.app" },' >> $(APPDMG_FILE)
	@echo '    { "x": 512, "y": 900, "type": "position", "path": ".VolumeIcon.icns" },' >> $(APPDMG_FILE)
	@echo '    { "x": 512, "y": 900, "type": "position", "path": ".background" }' >> $(APPDMG_FILE)
	@echo '  ]' >> $(APPDMG_FILE)
	@echo '}' >> $(APPDMG_FILE)

	@./node_modules/.bin/appdmg $(APPDMG_FILE) $(DMG_PATH)/CloudSQLProxyMenuBar-$(ARCH).dmg


.PHONY: build/dmg/amd64
build/dmg/amd64: ## Build amd64 dmg
	@$(MAKE) build/dmg ARCH=amd64


.PHONY: build/dmg/arm64
build/dmg/arm64: ## Build arm64 dmg
	@$(MAKE) build/dmg ARCH=arm64


.PHONY: build/dmg/all
build/dmg/all: ## Build all dmg
	@$(MAKE) build/dmg/amd64
	@$(MAKE) build/dmg/arm64


.PHONY: clean
clean: ## clean built outputs
	@rm -rf dist


# check variable definition
guard-%:
	@if [[ -z '${${*}}' ]]; then echo 'ERROR: variable $* not set' && exit 1; fi
