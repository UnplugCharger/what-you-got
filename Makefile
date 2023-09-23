DB_NAME=hackathon_db
DB_URI=postgresql://root:password@localhost:5432/$(DB_NAME)?sslmode=disable
all: test

postgres:
	docker run --name datapoint_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine

createdb:
	docker exec -it datapoint_db createdb --username=root --owner=root $(DB_NAME)

dropdb:
	docker exec -it datapoint_db dropdb $(DB_NAME)

migratedown:
	migrate -path db/migrations -database ${DB_URI} -verbose down 1

prodmigrateup:
	migrate -path db/migrations -database ${DB_URI} -verbose up

sqlc:
	sqlc generate

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

server:
	go run main.go