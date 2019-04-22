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