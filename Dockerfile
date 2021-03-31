FROM golang:1.16.0-alpine3.13
LABEL MAINTAINER=github.com/dezhiShen
RUN go env -w GO111MODULE=auto \
  && go env -w GOPROXY=https://goproxy.cn,direct 
RUN apk add -U --repository http://mirrors.ustc.edu.cn/alpine/v3.13/main/ tzdata
RUN apk add ca-certificates && update-ca-certificates
WORKDIR /build
COPY ./ .
RUN cd /build && go build -tags netgo -o miraigo cmd/main.go && cp /build/miraigo /usr/bin/miraigo && chmod +x /usr/bin/miraigo
WORKDIR /data
ENTRYPOINT ["/usr/bin/miraigo"]