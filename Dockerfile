FROM golang:1.19-alpine as build-stage

RUN apk --update upgrade \
    && apk --no-cache --no-progress add git mercurial bash gcc musl-dev curl tar ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on GOPROXY=https://goproxy.cn go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOGC=off  go build -ldflags="-X 'github.com/bitmyth/go-version.Version=v1.0.9'" -v -a -installsuffix nocgo -o dist/app ./

FROM scratch
COPY --from=build-stage /app/dist/app /
COPY --from=build-stage /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai
WORKDIR /
VOLUME /config
EXPOSE 8080
ENTRYPOINT ["/app"]
