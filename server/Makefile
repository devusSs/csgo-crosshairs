# Update the version to your needs via env / shell.
BUILD_VERSION=$(CHS_BUILD_VERSION)
BUILD_MODE=${CHS_BUILD_MODE}

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
BUILD_ARCH 				:=
ifeq ($(echo %PROCESSOR_ARCHITECTURE%), "AMD64")
	BUILD_ARCH = amd64
else
	UNAME_M := $(shell uname -m)
	ifeq ($(UNAME_M), x86_64)
		BUILD_ARCH = amd64
	endif
	ifeq ($(UNAME_M), arm64)
		BUILD_ARCH = arm64
	endif
endif

# DO NOT CHANGE.
build:
	@[ "${CHS_BUILD_VERSION}" ] || ( echo "CHS_BUILD_VERSION is not set"; exit 1 )
	@[ "${CHS_BUILD_MODE}" ] || ( echo "CHS_BUILD_MODE is not set"; exit 1 )
	@echo "Building app for Windows (AMD64), Linux (AMD64) & MacOS (ARM64)..."
	@go mod tidy
	@GOOS=windows GOARCH=amd64 go build -ldflags="-X 'github.com/devusSs/crosshairs/updater.BuildVersion=$(BUILD_VERSION)' -X 'github.com/devusSs/crosshairs/updater.BuildDate=${shell date}' -X 'github.com/devusSs/crosshairs/updater.BuildMode=$(BUILD_MODE)'" -o release/chs_windows_amd64/ ./...
	@GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/devusSs/crosshairs/updater.BuildVersion=$(BUILD_VERSION)' -X 'github.com/devusSs/crosshairs/updater.BuildDate=${shell date}' -X 'github.com/devusSs/crosshairs/updater.BuildMode=$(BUILD_MODE)'" -o release/chs_linux_amd64/ ./...
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'github.com/devusSs/crosshairs/updater.BuildVersion=$(BUILD_VERSION)' -X 'github.com/devusSs/crosshairs/updater.BuildDate=${shell date}' -X 'github.com/devusSs/crosshairs/updater.BuildMode=$(BUILD_MODE)'" -o release/chs_darwin_arm64/ ./...
	@echo "Done building app"

# DO NOT CHANGE.
build-release:
	@[ "${CHS_BUILD_VERSION}" ] || ( echo "CHS_BUILD_VERSION is not set"; exit 1 )
	@echo "Building release-ready app for Windows (AMD64), Linux (AMD64) & MacOS (ARM64)..."
	@go mod tidy
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -trimpath -ldflags="-s -w -X 'github.com/devusSs/crosshairs/updater.BuildVersion=$(BUILD_VERSION)' -X 'github.com/devusSs/crosshairs/updater.BuildDate=${shell date}' -X 'github.com/devusSs/crosshairs/updater.BuildMode=release'" -o release/chs_windows_amd64/ ./...
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -trimpath -ldflags="-s -w -X 'github.com/devusSs/crosshairs/updater.BuildVersion=$(BUILD_VERSION)' -X 'github.com/devusSs/crosshairs/updater.BuildDate=${shell date}' -X 'github.com/devusSs/crosshairs/updater.BuildMode=release'" -o release/chs_linux_amd64/ ./...
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -v -trimpath -ldflags="-s -w -X 'github.com/devusSs/crosshairs/updater.BuildVersion=$(BUILD_VERSION)' -X 'github.com/devusSs/crosshairs/updater.BuildDate=${shell date}' -X 'github.com/devusSs/crosshairs/updater.BuildMode=release'" -o release/chs_darwin_arm64/ ./...
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
	@cd ./testing && ./crosshairs -c "./files/config.dev.json" -d

# DO NOT CHANGE.
full-dev: drop-postgres dev-postgres drop-redis dev-redis drop-minio dev-minio build
	@clear
	@rm -rf ./testing
	@mkdir ./testing
	@mkdir ./testing/files
	@cp -R ./templates ./testing
	@cp -R ./files ./testing
	@cp ./release/chs_$(BUILD_OS)_$(BUILD_ARCH)/crosshairs ./testing
	@cd ./testing && ./crosshairs -c "./files/config.dev.json" -d

# DO NOT CHANGE.
docker-up:
	@clear
	@echo "Checking for docker.env file..."
	@[ -f ./docker.env ] && echo "Found docker.env file, proceeding..." || echo "ERROR: missing docker.env file. Please create a docker.env file and edit it" && exit 0
	@docker compose --env-file=docker.env up --build -d
	@echo ""
	@echo "Please use 'make docker-down' to shutdown the app and containers."
	@echo "You may also run 'make docker-rebuild' to rebuild the app and restart the containers"

# DO NOT CHANGE.
docker-down:
	@echo "This command will NOT delete Docker images or volumes (data will be persistent)."
	@echo "If you need functionality to remove images and volumes please use the Docker GUI / CLI if available."
	@docker compose --env-file=docker.env down

# DO NOT CHANGE.
docker-rebuild: docker-down docker-up

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
	@rm -rf ./minio/

# DO NOT CHANGE
full-clean: clean drop-postgres drop-redis drop-minio

# DO NOT CHANGE.
dev-postgres:
	@docker run -d --rm --name pg-chs -e POSTGRES_PASSWORD=crosshairs -e POSTGRES_USER=crosshairs -e POSTGRES_DB=crosshairs -p 5432:5432 -d postgres

# DO NOT CHANGE.
drop-postgres:
	@docker stop pg-chs || true

# DO NOT CHANGE
dev-redis:
	@docker volume create redis-data
	@docker run -d --rm --name redis-chs -v redis-data:/data -p 6379:6379 redis redis-server --requirepass crosshairs

# DO NOT CHANGE
drop-redis:
	@docker stop redis-chs || true
	@docker volume rm redis-data || true

# DO NOT CHANGE
dev-minio:
	@docker run -d --rm -p 9000:9000 --name minio-chs -v ./minio/data:/data -e "MINIO_BROWSER=off" -e "MINIO_ROOT_USER=crosshairs" -e "MINIO_ROOT_PASSWORD=crosshairs" \quay.io/minio/minio server /data

# DO NOT CHANGE
drop-minio:
	@docker stop minio-chs || true

# DO NOT CHANGE.
lint:
	@golangci-lint run
