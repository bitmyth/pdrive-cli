# Personal Drive on the Cloud

### migrate database

```shell
go run db/migrate/migrate.go
```

### Run by docker

```shell
docker run -d --name pdrive --restart always -p8008:80 -v$(pwd)/upload:/upload -v $(pwd)/config:/config bitmyth/pdrive
```

### Your first token

your first token will be printed out to console

## CLI

### run cli from source

```shell
 go run cli/cmd/pd/main.go file upload --file  /Users/gsh/Downloads/CredentiaReport.csv^
```

### build

go build -o pd cli/cmd/pd/main.go

go install cli/cmd/pd/main.go

### login command

```shell
pd auth login
```

### upload file command

```shell
pd file upload --file
```

### MySQL

docker run -d --net app --name mysql -p3306:3306 -e MYSQL_ROOT_PASSWORD=password -e TZ=Asia/Shanghai mysql:5.7

### Create database
mysql> create database test default charset utf8mb4;
# Run

### Docker

docker run -d --name pdrive --restart always -p8008:80 -v$(pwd)/upload:/upload -v $(pwd)/config:/config bitmyth/pdrive