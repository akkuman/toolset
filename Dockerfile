FROM alpine/curl AS downloader

WORKDIR /app
RUN curl -s -o supervisord.tar.gz https://ghproxy.com/https://github.com/ochinchina/supervisord/releases/download/v0.7.3/supervisord_0.7.3_Linux_64-bit.tar.gz && \
    tar -zxvf supervisord.tar.gz -C . && \
    mv /app/supervisord_0.7.3_Linux_64-bit/supervisord_static /app/supervisord



FROM golang:1.18 as builder

WORKDIR /src

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go install mvdan.cc/garble@latest

COPY ./go.mod ./go.sum /src/
RUN go get -u github.com/swaggo/swag/cmd/swag && go mod download
ADD . /src/
RUN rm -rf /src/data
RUN swag init && go build -o /app/toolset . && cp /go/bin/garble /app/garble


FROM golang:1.18 AS finally

WORKDIR /app

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

RUN apt update && \
    apt install -y mingw-w64 && \
    # Cleaning cache:
    apt-get clean -y && rm -rf /var/lib/apt/lists/*

COPY --from=downloader /app/supervisord /app/supervisord
ADD data /app/data
COPY --from=builder /app/toolset /app/
COPY --from=builder /app/garble /go/bin/

COPY supervisord.conf /app/

EXPOSE 8080 9001
CMD ["/app/supervisord", "-c", "/app/supervisord.conf"]
