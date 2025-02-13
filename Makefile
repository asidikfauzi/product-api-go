#Go Params
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod

CMDDIR=cmd
BIN_NAME=app

DBDIR=migration
BIN_MIGRATION=migration

MIGDIR=migrate
BIN_MIGRATE=migrate

SEEDDIR=seed
BIN_SEED=seed

DOCKER_FILE=build/package/docker

all: mod build-main run
build: build-main
reload: build-main run
migrate: migrate-build migrate-run
rollback: rollback-build rollback-run
seed: mod seed-run run-seed
docker-build: docker-build
docker-down: docker-down

mod:
	$(GOMOD) vendor -v

build-main:
	rm -f $(CMDDIR)/$(BIN_NAME)/$(BIN_NAME)
	$(GOBUILD) -o $(CMDDIR)/$(BIN_NAME)/$(BIN_NAME) $(CMDDIR)/${BIN_NAME}/app.go

run:
	./$(CMDDIR)/$(BIN_NAME)/$(BIN_NAME)

new-migration:
	rm -f $(CMDDIR)/$(DBDIR)/$(BIN_MIGRATION)
	$(GOBUILD) -o $(CMDDIR)/$(DBDIR)/$(BIN_MIGRATION) $(CMDDIR)/$(DBDIR)/migration.go
	./$(CMDDIR)/$(DBDIR)/$(BIN_MIGRATION) -table=$(table)
	rm -f $(CMDDIR)/$(DBDIR)/$(BIN_MIGRATION)

migrate-build:
	$(GOBUILD) -o $(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE) $(CMDDIR)/$(MIGDIR)/migrate.go
	./$(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE) -mode=up
	rm -f $(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE)

migrate-run:
	$(GOBUILD) -o $(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE) $(CMDDIR)/$(MIGDIR)/migrate.go
	./$(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE) -mode=exec-up
	rm -f $(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE)

rollback-build:
	$(GOBUILD) -o $(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE) $(CMDDIR)/$(MIGDIR)/migrate.go
	./$(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE) -mode=down -step=$(step)
	rm -f $(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE)

rollback-run:
	$(GOBUILD) -o $(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE) $(CMDDIR)/$(MIGDIR)/migrate.go
	./$(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE) -mode=exec-down
	rm -f $(CMDDIR)/$(MIGDIR)/$(BIN_MIGRATE)

seed-run:
	rm -f $(CMDDIR)/$(SEEDDIR)/$(BIN_SEED)
	$(GOBUILD) -o $(CMDDIR)/$(SEEDDIR)/$(BIN_SEED) $(CMDDIR)/$(SEEDDIR)/seed.go

run-seed:
	./$(CMDDIR)/$(SEEDDIR)/$(BIN_SEED)
	rm -f $(CMDDIR)/$(SEEDDIR)/$(BIN_SEED)

docker-build:
	(cd $(DOCKER_FILE) && docker compose -f docker-compose.yml -f dev.docker-compose.yml --env-file ../../../.env up --build)

docker-down:
	(cd $(DOCKER_FILE) && docker compose down)

