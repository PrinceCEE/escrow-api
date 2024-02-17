.PHONY: run/api
run/app/dev:
	go run ./cmd/app -env=development -loglevel=debug