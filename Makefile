# Makefile

NAME  := memoria-server
ARCH  ?= amd64
OS    ?= linux
UNAME := $(shell uname -s)

VERSION := 0.0.1
DIST := dist

CGO_ENABLED=0
GOARCH="${ARCH}"
GOOS="${OS}"

# Initial go mod
# go mod init

# Run
.PHONY: run
run:
	go run src/main.go

#  Vendor
.PHONY: vendor
vendor:
	go mod vendor

#  Build
.PHONY: build
build: clean vendor
	ARCH=$(ARCH) OS=$(OS) VERSION=$(VERSION) go build -o $(DIST)/$(NAME) ./src/... 

#  Test
.PHONY: test
test:
	go -v test ./...

# Clean
.PHONY: clean
clean:
	rm -Rf $(DIST)

# ----------

# Check the health endpoint.
.PHONY: check-health
check-health:
	curl localhost:8080/health

# Check the health endpoint.
.PHONY: check-signup
check-signup:
	curl -X POST localhost:8080/signup -d '{ "name": "test", "email": "test", "password": "test" }'

# Check the health endpoint.
.PHONY: check-login
check-login:
	curl localhost:8080/login -d '{ "email": "test", "password": "test" }'