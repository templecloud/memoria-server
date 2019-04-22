# Makefile

# Initial go mod
# go mod init

# Directly run main.
.PHONY: run-main
run-main:
	go run src/main.go

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