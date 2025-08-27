API_CMD=./cmd/api/main.go
SEEDER_CMD=./cmd/seeder/main.go
MIGRATE_CMD=./cmd/migrate/main.go
API_BIN=./bin/api
SEEDER_BIN=./bin/seeder
MIGRATE_BIN=./bin/migrate

build:
	@go build -o $(API_BIN) $(API_CMD)

dev:
	@JWT_SECRET="DBIUAB!!312312ADHI"
	@air \
	 --build.cmd "clear && make -s build" \
	 --build.bin "$(API_BIN)" \
	 --build.exclude_dir "bin,tmp" \
	 --build.include_ext "go,yaml"

test:
	@go test ./tests

seed: 
	@go build -o $(SEEDER_BIN) $(SEEDER_CMD)
	@$(SEEDER_BIN)

migrate-up:
	@go build -o $(MIGRATE_BIN) $(MIGRATE_CMD)
	@$(MIGRATE_BIN) up

migrate-down: 
	@go build -o $(MIGRATE_BIN) $(MIGRATE_CMD)
	@$(MIGRATE_BIN) down

docker-run-api:
	docker compose up

docker-run-migrations:
	docker compose run --rm api /app/migrate
