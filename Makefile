MIGRATION_ROOT?=./db/migration

# optionally pass this from command line
DB?=cesc_blog

postgres:
	docker run --name postgres15_0 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:15.0-alpine

createdb:
	docker exec -it postgres15_ca_ES createdb --username=postgres --owner=postgres ${DB}

dropdb:
	docker exec -it postgres15_0 dropdb --username=postgres ${DB}

migrateup:
	migrate -verbose -path $(MIGRATION_ROOT) -database "postgres://postgres:postgres@localhost:5432/$(DB)?sslmode=disable" up

migratedown:
	migrate -verbose -path $(MIGRATION_ROOT) -database "postgres://postgres:postgres@localhost:5432/$(DB)?sslmode=disable" down

sqlc-init:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc init

sqlc-generate:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go blogapi/db/sqlc Store

test:
	go test -v -cover -short ./...

server:
	go run main.go

.PHONY: createdb dropdb postgres migrateup migratedown sqlc-generate sqlc-init mock test server

