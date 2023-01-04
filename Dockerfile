# Build Stage
FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN apk add --virtual --update --no-cache protoc curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN mkdir -p pb
RUN rm -f pb/*.go \
	rm -f docs/swagger/*.swagger.json
RUN protoc --proto_path=proto	--grpc-gateway_out=pb	--grpc-gateway_opt=paths=source_relative	--go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --openapiv2_out=docs/swagger	--openapiv2_opt=allow_merge=true,merge_file_name=bank \
    proto/*.proto

RUN go build -o main main.go

#Run Stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate 
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration
RUN chmod +x start.sh
RUN chmod +x wait-for.sh
EXPOSE 8080
# CMD will act as extra parameters passed into the script
# eg : ENTRYPOINT [ "/app/start.sh", "/app/main" ]
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]