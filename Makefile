
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

postgres:
	docker pull postgres && docker run --network bank-network --name bank-postgres -p 5430:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it bank-postgres createdb --username=root --owner=root bank
sqlc:
	sqlc generate

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5430/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5430/bank?sslmode=disable" -verbose down

dropdb:
	docker exec -it bank-postgres dropdb --username=root bank

clean_test:
	go clean -testcache && go test -v -cover ./...

test:
	go test -v -cover ./...

mock:
	mockgen -destination db/mock/store.go -package mockdb github/dutt23/bank/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto	--grpc-gateway_out=pb	--grpc-gateway_opt=paths=source_relative	--go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto
evans:
	evans --host localhost --port 9090 -r repl

server: 
	go	run	main.go 
