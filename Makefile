include .env

.PHONY: run/api/dev
run/api/dev:
	@echo "starting the web server"
	go run ./cmd/app -env=development -loglevel=debug

.PHONY: migration/create
migration/create:
	@echo "creating new migration files for ${name}"
	migrate create -ext sql -dir ./migrations -seq -digits 8 ${name}

.PHONY: migration/up
migration/up:
	@echo "running migration"
	migrate -path ./migrations -database ${DSN} up

.PHONY: migration/down
migration/down:
	@echo "running migration"
	migrate -path ./migrations -database ${DSN} down