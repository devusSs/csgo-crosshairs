# Backend for chs.devuscs.me

## Setup

Make sure you have a running [Postgres](https://www.postgresql.org/) instance and have the latest version of [Go(lang)](https://go.dev) installed on your system.

You may then setup the config file according to the [example config](files/config.json)'s specifications. The program will automatically error and exit in case something does not work properly.

## Building and testing

It is highly recommended to use the [Makefile](Makefile) / [GNU Make](https://www.gnu.org/software/make/) to build and run the backend:

- `make build` to build binaries for Windows (AMD64), Linux (AMD64), MacOS (ARM64).

- `make version` to run make build and run the app only printing build information.

- `make dev` to run make build and run the actual app.

- `make clean` to clean the directory of build artifacts and logs.

- `make dev-postgres` to generate a working Postgres instance using Docker (needs to be installed).

- `make lint` linting Go code using golangci-lint (needs to be installed).

## API routes, requests & responses structure

The documentation can be found in the [docs directory](docs).
