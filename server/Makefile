# Update the version to your needs via env / shell.
BUILD_VERSION=$(CHS_BUILD_VERSION)
BUILD_MODE=${CHS_BUILD_MODE}

# DO NOT CHANGE.
BUILD_DATE=$$(date +%Y.%m.%d-%H:%M:%S)

# DO NOT CHANGE.
BUILD_OS 				:=
ifeq ($(OS),Windows_NT)
	BUILD_OS = windows
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		BUILD_OS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		BUILD_OS = darwin
	endif
endif

# DO NOT CHANGE.
BUILD_ARCH				:=
ifeq ($(OS),Windows_NT)
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		BUILD_ARCH = amd64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		BUILD_ARCH = 386
	endif
else
		UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		BUILD_ARCH = amd64
	endif
		ifneq ($(filter %86,$(UNAME_P)),)
	BUILD_ARCH = 386
		endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		BUILD_ARCH = arm64
	endif
endif

# DO NOT CHANGE.
build: 
	@[ "${CHS_BUILD_VERSION}" ] || ( echo "CHS_BUILD_VERSION is not set"; exit 1 )
	@[ "${CHS_BUILD_MODE}" ] || ( echo "CHS_BUILD_MODE is not set"; exit 1 )
	@echo "Building app for Windows (AMD64), Linux (AMD64) & MacOS (ARM64)..."
	@go mod tidy
	@GOOS=windows GOARCH=amd64 go build -v -trimpath -ldflags="-s -w -X github.com/devusSs/crosshairs/updater.buildVersion=$(BUILD_VERSION) -X github.com/devusSs/crosshairs/updater.buildDate=$(BUILD_DATE) -X github.com/devusSs/crosshairs/updater.BuildMode=$(BUILD_MODE)" -o release/chs_windows_amd64/ ./...
	@GOOS=linux GOARCH=amd64 go build -v -trimpath -ldflags="-s -w -X github.com/devusSs/crosshairs/updater.buildVersion=$(BUILD_VERSION) -X github.com/devusSs/crosshairs/updater.buildDate=$(BUILD_DATE) -X github.com/devusSs/crosshairs/updater.BuildMode=$(BUILD_MODE)" -o release/chs_linux_amd64/ ./...
	@GOOS=darwin GOARCH=arm64 go build -v -trimpath -ldflags="-s -w -X github.com/devusSs/crosshairs/updater.buildVersion=$(BUILD_VERSION) -X github.com/devusSs/crosshairs/updater.buildDate=$(BUILD_DATE) -X github.com/devusSs/crosshairs/updater.BuildMode=$(BUILD_MODE)" -o release/chs_darwin_arm64/ ./...
	@echo "Done building app"

# DO NOT CHANGE.
dev: build
	@clear
	@rm -rf ./testing
	@mkdir ./testing
	@mkdir ./testing/files
	@cp -R ./templates ./testing
	@cp -R ./files ./testing
	@cp ./release/chs_$(BUILD_OS)_$(BUILD_ARCH)/crosshairs ./testing
	@cd ./testing && ./crosshairs -c "./files/config.dev.json"

# DO NOT CHANGE.
full-dev: drop-postgres dev-postgres drop-redis dev-redis build
	@clear
	@rm -rf ./testing
	@mkdir ./testing
	@mkdir ./testing/files
	@cp -R ./templates ./testing
	@cp -R ./files ./testing
	@cp ./release/chs_$(BUILD_OS)_$(BUILD_ARCH)/crosshairs ./testing
	@cd ./testing && ./crosshairs -c "./files/config.dev.json"

# DO NOT CHANGE.
version: build
	@clear
	@rm -rf ./testing
	@mkdir ./testing
	@cp ./release/chs_$(BUILD_OS)_$(BUILD_ARCH)/crosshairs ./testing
	@cd ./testing && ./crosshairs -v

# DO NOT CHANGE.
secret-keys: build
	@clear
	@rm -rf ./testing
	@mkdir ./testing
	@cp ./release/chs_$(BUILD_OS)_$(BUILD_ARCH)/crosshairs ./testing
	@cd ./testing && ./crosshairs -sc

# DO NOT CHANGE.
clean:
	@clear
	@go mod tidy
	@go fmt ./...
	@rm -rf ./release/
	@rm -rf ./dist/
	@rm -rf ./testing/

# DO NOT CHANGE.
dev-postgres:
	@docker run -d --rm --name pg-chs -e POSTGRES_PASSWORD=crosshairs -e POSTGRES_USER=crosshairs -e POSTGRES_DB=crosshairs -p 5432:5432 -d postgres

# DO NOT CHANGE.
drop-postgres:
	@docker stop pg-chs

# DO NOT CHANGE
dev-redis:
	@docker volume create redis-data
	@docker run -d --rm --name redis-chs -e REDIS_PASSWORD=crosshairs -v redis-data:/data -p 6379:6379 redis:5.0.5-alpine3.9 /bin/sh -c redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}

# DO NOT CHANGE
drop-redis:
	@docker stop redis-chs
	@docker volume rm redis-data

# DO NOT CHANGE.
lint:
	@golangci-lint run -v