postgres:
	docker run --name postgres14_1 -p 5433:5433 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=cirey.anaid.2107 -d postgres:14.1-alpine

createdb:
	docker exec -it postgres14_1 createdb --username=root --owner=root yeric-blog

dropdb:
	docker exec -it postgres14_1 dropdb yeric-blog

migrateup:
	migrate -path models/db/migration -database "postgresql://root:cirey.anaid.2107@localhost:5433/yeric-blog?sslmode=disable" up

migratedown:
	migrate -path models/db/migration -database "postgresql://root:cirey.anaid.2107@localhost:5433/yeric-blog?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown