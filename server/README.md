# Backend for dropawp.com

## Development

It is highly recommended to use Docker (with docker-compose) because it will take care of everything database and storage related.<br/>
Just download the latest version of [Docker](https://www.docker.com/products/docker-desktop/) for your operating system, head to the server's directory and run:

```bash
make docker-up
```

To shutdown the containers run:

```bash
make docker-down
```

Please make sure to follow the instructions on screen, especially related to the `docker.env` file.

## Building and testing (not recommended for development)

### Setup

Make sure you have a running [Postgres](https://www.postgresql.org/) instance, [Redis instance](https://redis.io/docs/getting-started/) and [Minio instance](https://min.io/download#/windows) and have the latest version of [Go(lang)](https://go.dev) installed on your system.

You may then setup the config file according to the [example config](files/config.json)'s specifications. The program will automatically error and exit in case something does not work properly.

It is highly recommended to use the [Makefile](Makefile) / [GNU Make](https://www.gnu.org/software/make/) to build and run the backend:

- `make version` to run make build and run the app only printing build information.

- `make full-dev` to drop the Postgres and Redis dbs, create new ones, run make build and run the actual app.

- `make full-clean` to clean the directory of build artifacts and logs as well as dropping every related service.

- `make lint` linting Go code using golangci-lint (needs to be installed).

- `make secret-keys` to generate the session token and admin token.

## API routes, requests & responses structure

The documentation can be found in the [docs directory](api/docs).
