TAG=bitmyth/pdrive-cli
CLI=cli/cmd/pd/main.go

image:
	docker build -t $(TAG) .

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOGC=off  go build -ldflags="-X 'github.com/bitmyth/go-version.Version=v1.0.9'" -o pd $(CLI)

pdcli:
	#go build -ldflags="-X github.com/bitmyth/go-version.Version=$(GITTAG)" -o pd $(CLI)
	go generate cli/build/version.go
	go build -o pd $(CLI)

clirelease:
	go generate cli/build/version.go
	GOOS=windows GOARCH=amd64 GOGC=off go build -o pd_windows_amd64 cli/cmd/pd/main.go
	GOOS=darwin GOARCH=arm64 GOGC=off go build -o pd_darwin_arm64 cli/cmd/pd/main.go
	GOOS=darwin GOARCH=amd64 GOGC=off go build -o pd_darwin_amd64 cli/cmd/pd/main.go
	GOOS=linux GOARCH=amd64 GOGC=off go build -o pd_linux_amd64 cli/cmd/pd/main.go
install:
	go generate cli/build/version.go
	go install ./cli/cmd/pd