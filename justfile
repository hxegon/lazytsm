_default:
    @just --choose

run:
    go run cmd/lazyproj.go

build:
    go build cmd/lazyproj.go

tidy:
    go mod tidy
