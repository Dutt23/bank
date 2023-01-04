List of dependencies/tools used:

1) sqlc
2) gin
3) viper
4) docker
5) gomock ( add go binary to path like this, PATH=$PATH:~/go/bin)
6) https://github.com/golang-migrate/migrate
7) https://formulae.brew.sh/formula/golang-migrate


Docker parts
Both of the containers run on the default bridge network, but they can't call each other via names.
Creating a custom network for both of them. docker network --help


#Migrate
migrate create -ext sql -dir db/migration -seq <migration_name>

#Gprpc
install . brew install ptotobufjs. protoc --version
https://github.com/grpc-ecosystem/grpc-gateway/blob/main/examples/internal/proto/examplepb/a_bit_of_everything.proto
<!-- Gateways -->
https://github.com/grpc-ecosystem/grpc-gateway
protoc-gen-grpc-gateway -help
*.gw.go File for the httpgateway
<!-- Testing purposes -->
https://grpc.io/docs/languages/go/quickstart/
https://github.com/ktr0731/evans#macos
evans --host localhost --port 9090 -r repl
