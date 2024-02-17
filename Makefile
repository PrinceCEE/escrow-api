.PHONY: run/app/dev
run/app/dev:
	go run ./cmd/app -env=development -loglevel=debug