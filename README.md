# Convoluted Movie Finder

An intentionally convulated movie finding app built using the TMDB API to play around with Go's concurrency features

## Pre-requisites

- [Go 1.17](https://golang.org/)
- [Protocol buffer compiler](https://grpc.io/docs/protoc-installation/)

## Configuration

This app uses [configo](https://github.com/affanshahid/configo) for configurations. Configure different parameters including the required `tmdb_api_key` using the config folder or environment variables.

## Running

```sh
go run ./cmd/movie-finder # to run the server

go run ./cmd/client # to run the client with some default search options

go run ./cmd/client -help # to see client usage
```

## Generating mocks and gRPC code

```sh
go generate ./...
```

## Testing

```sh
APP_ENV=test go test ./...
```
