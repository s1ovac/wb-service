.PHONY: build
build: 
		go build -v order-publish/cmd/main.go

.DEFAULT_GOAL := build 