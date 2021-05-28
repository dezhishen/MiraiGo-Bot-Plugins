FROM golang:1.16.0-alpine3.13 AS builder
LABEL MAINTAINER=github.com/dezhiShen
RUN go env -w GO111MODULE=auto \
  && go env -w GOPROXY=https://goproxy.cn,direct 
RUN apk add -U --repository http://mirrors.ustc.edu.cn/alpine/v3.13/main/ tzdata
RUN apk add ca-certificates && update-ca-certificates
WORKDIR /build
COPY ./ .
RUN cd /build && go build -tags netgo -ldflags="-w -s" -o miraigo cmd/main.go 

FROM alpine:latest

WORKDIR /data

COPY --from=builder /build/miraigo /usr/bin/miraigo 
RUN chmod +x /usr/bin/miraigo
VOLUME /data
HEALTHCHECK  --interval=5s --timeout=1s --start-period=5s --retries=3 CMD cat /data/health
ENTRYPOINT ["/usr/bin/miraigo"]