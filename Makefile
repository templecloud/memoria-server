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
.PHONY: dev
dev:
	go run cmd/memoria/main.go

#  Vendor
.PHONY: vendor
vendor:
	go mod vendor

#  Build
.PHONY: build
build: clean vendor
	ARCH=$(ARCH) OS=$(OS) VERSION=$(VERSION) go build -o $(DIST)/$(NAME) ./cmd/... 

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
	curl -v --cookie "token=${JWT}" localhost:8080/api/v1/health

# Check the signup endpoint.
.PHONY: check-signup
check-signup:
	curl -v -X POST localhost:8080/api/v1/signup -d '{ "name": "test", "email": "test", "password": "test" }'

# Check the login endpoint.
.PHONY: check-login
check-login:
	curl -v localhost:8080/api/v1/login -d '{ "email": "test", "password": "test" }'