DB_URL=postgresql://root:secret@localhost:5430/bank?sslmode=disable

.PHONY: postgres
.PHONY: createdb
.PHONY: dropdb
.PHONY: sqlc
.PHONY: migrateup
.PHONY: migratedown
.PHONY: test
.PHONY: clean_test
.PHONY: server
.PHONY:	mock
.PHONY:	proto
.PHONY: evans
.PHONY: redis
.PHONY: db_schema
.PHONY: db_docs
.PHONY: new_migration

postgres:
	docker pull postgres && docker run --network bank-network --name bank-postgres -p 5430:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it bank-postgres createdb --username=root --owner=root bank
sqlc:
	sqlc generate

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

dropdb:
	docker exec -it bank-postgres dropdb --username=root bank

clean_test:
	go clean -testcache && go test -v -cover ./...

test:
	go test -v -cover -short ./...

mock:
	mockgen -destination db/mock/store.go -package mockdb github/dutt23/bank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

new_migration: 
	migrate create -ext sql -dir db/migration -seq $(name)

proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto	--grpc-gateway_out=pb	--grpc-gateway_opt=paths=source_relative	--go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--openapiv2_out=docs/swagger	--openapiv2_opt=allow_merge=true,merge_file_name=bank\
    proto/*.proto 
	statik -src=./docs/swagger -dest=./docs

evans:
	evans --host localhost --port 9090 -r repl

server: 
	go	run	main.go 

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
