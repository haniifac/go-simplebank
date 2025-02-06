postgres:
	docker run --name pg12-bank -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:12-alpine

createdb:
	docker exec -it pg12-bank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it pg12-bank dropdb simple_bank

migrateup:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test