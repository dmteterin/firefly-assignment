.PHONY: run
run:
	@go run cmd/*.go

.PHONY: build
build:
	@go build -o ./app cmd/*.go