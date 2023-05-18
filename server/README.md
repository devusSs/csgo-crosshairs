# Backend for chs.devuscs.me

## Setup

Make sure you have a running [Postgres](https://www.postgresql.org/) instance and have the latest version of [Go(lang)](https://go.dev) installed on your system.

You may then setup the config file according to the [example config](files/config.json)'s specifications. The program will automatically error and exit in case something does not work properly.

## Building and testing

It is highly recommended to use the [Makefile](Makefile) / [GNU Make](https://www.gnu.org/software/make/) to build and run the backend:

- `make build` to build binaries for Windows (AMD64), Linux (AMD64), MacOS (ARM64).

- `make version` to run make build and run the app only printing build information.

- `make dev` to run make build and run the actual app.

- `make full-dev` to drop the Postgres and Redis dbs, create new ones, run make build and run the actual app.

- `make clean` to clean the directory of build artifacts and logs.

- `make dev-postgres` to generate a working Postgres instance using Docker (needs to be installed).

- `make drop-postgres` to drop the created Postgres database.

- `make dev-redis` to generate a working Redis instance using Docker (needs to be installed)

- `make drop-redis` to drop the created Redis database.

- `make lint` linting Go code using golangci-lint (needs to be installed).

- `make secret-keys` to generate the session token and admin token.

## API routes, requests & responses structure

The documentation can be found in the [docs directory](api/docs).
