build:
	go build -o bin/
postgres:
	docker run --name postgres14_1 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234567 -d postgres:14.1-alpine3.15

createdb:
	docker exec -it postgres14_1 createdb --username=root --owner=root yeric-blog

dropdb:
	docker exec -it postgres14_1 dropdb yeric-blog

migrateup:
	migrate -path models/db/migration -database "postgresql://root:1234567@localhost:5433/yeric-blog?sslmode=disable" -verbose up

migratedown:
	migrate -path models/db/migration -database "postgresql://root:1234567@localhost:5433/yeric-blog?sslmode=disable" -verbose down

.PHONY: build postgres createdb dropdb migrateup migratedown