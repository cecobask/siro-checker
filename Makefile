.PHONY: *

check:
	@go run main.go

lint:
	@golangci-lint run

lint-fix:
	@golangci-lint run --fix