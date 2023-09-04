# Personal Drive CLI
Current tag v1.1.4

# Install 
go install github.com/bitmyth/pdrive-cli/cli/cmd/pd

### run cli from source

```shell
 go run cli/cmd/pd/main.go file upload --file  $FILEPATH
```

### install from source
go install ./cli/cmd/pd

### login command

```shell
pd auth login
```

### upload file command

```shell
pd file upload --file
```