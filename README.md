# Convoluted Movie Finder

An intentionally convulated movie finding app built using the TMDB API to play around with Go's concurrency features

## Pre-requisites

- [Go 1.17](https://golang.org/)
- [Protocol buffer compiler](https://grpc.io/docs/protoc-installation/)

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
