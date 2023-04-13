MIGRATION_ROOT?=./db/migration

# optionally pass this from command line
DB?=cesc_blog

RELEASE_VERSION?=0.0.1

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

build:
	go build -o build ./...

test:
	go test -v -cover -short ./...

server:
	go run cmd/main.go

docker-build:
	docker build -t cesc-blog-api:dev-0.0.1 .

docker-run:
	docker run -it -d -p 5000:5000 --rm --name cesc-blog-api cesc-blog-api:dev-${RELEASE_VERSION}

docker-logs:
	docker logs -f cesc-blog-api

docker-stop:
	docker stop cesc-blog-api

docker-compose-up:
	docker-compose -f .\docker-compose-dev.yml up                                                                                                                                                                                         in pwsh at 16:40:50 

.PHONY: build-development
build-development: ## Build the development docker image.
	docker compose -f docker/development/docker-compose.yml build


.PHONY: createdb dropdb postgres migrateup migratedown sqlc-generate sqlc-init mock test server build docker-build docker-run docker-logs docker-stop docker-compose-up

