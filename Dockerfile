# Build Stage
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

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