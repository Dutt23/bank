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