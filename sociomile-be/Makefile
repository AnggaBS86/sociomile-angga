APP_NAME=sociomile-be
MIGRATIONS_DIR=./migrations

DB_DSN ?= root:root@tcp(localhost:3306)/sociomile?parseTime=true
MIGRATE_DATABASE_URL ?= mysql://root:root@tcp(localhost:3306)/sociomile

# Always use mysql build tag to ensure mysql driver is included.
MIGRATE_CMD ?= GOTOOLCHAIN=local go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

AIR_BIN := $(shell command -v air 2>/dev/null)
ifeq ($(AIR_BIN),)
AIR_CMD := GOTOOLCHAIN=local go run github.com/air-verse/air@latest -c .air.toml
else
AIR_CMD := air -c .air.toml
endif

.PHONY: run worker dev test migrate-up migrate-down migrate-force migrate-create migrate-install air-install

run:
	GOTOOLCHAIN=local go run ./cmd/api

worker:
	GOTOOLCHAIN=local go run ./cmd/worker

dev:
	$(AIR_CMD)

test:
	GOTOOLCHAIN=local go test ./...

migrate-up:
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database '$(MIGRATE_DATABASE_URL)' up

migrate-down:
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database '$(MIGRATE_DATABASE_URL)' down 1

migrate-force:
	@test -n "$(VERSION)" || (echo "VERSION is required. usage: make migrate-force VERSION=1" && exit 1)
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database '$(MIGRATE_DATABASE_URL)' force $(VERSION)

migrate-create:
	@test -n "$(NAME)" || (echo "NAME is required. usage: make migrate-create NAME=add_indexes" && exit 1)
	$(MIGRATE_CMD) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

migrate-install:
	GOTOOLCHAIN=local go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

air-install:
	GOTOOLCHAIN=local go install github.com/air-verse/air@latest
