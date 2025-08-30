API_CMD=./cmd/api/main.go
SEEDER_CMD=./cmd/seeder/main.go
MIGRATE_CMD=./cmd/migrate/main.go
API_BIN=./bin/api
SEEDER_BIN=./bin/seeder
MIGRATE_BIN=./bin/migrate

dev:
	@docker compose up -d db
	@bash -c 'trap "docker compose down; exit" INT; \
	  air --build.cmd "clear && go build -tags=dev -o ./tmp/dev $(API_CMD)" \
	  --build.bin "./tmp/dev" \
	  --build.exclude_dir "bin,tmp" \
	  --build.include_ext "go,yaml"'

build:
	@go build -o $(API_BIN) $(API_CMD)

test:
	@docker compose up -d db
	@bash -c 'trap "docker compose down; exit" INT; go test -tags=dev ./tests'

seed:
	@go build -tags=dev -o $(SEEDER_BIN) $(SEEDER_CMD)
	@$(SEEDER_BIN)

migrate-up:
	@docker compose up -d db
	@go build -tags=dev -o $(MIGRATE_BIN) $(MIGRATE_CMD)
	@$(MIGRATE_BIN) up

migrate-down:
	@docker compose up -d db
	@go build -tags=dev -o $(MIGRATE_BIN) $(MIGRATE_CMD)
	@$(MIGRATE_BIN) down

run-api:
	docker compose up

run-migrations:
	docker compose run --rm api /app/migrate
