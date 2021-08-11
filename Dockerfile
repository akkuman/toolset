FROM golang:1.16 as builder

WORKDIR /src

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    GO111MODULE=on go get mvdan.cc/garble

COPY ./go.mod ./go.sum /src/
RUN go get -u github.com/swaggo/swag/cmd/swag && go mod download
ADD . /src/
RUN swag init && go mod tidy && go build -o /app/toolset . && cp /go/bin/garble /app/garble


FROM golang:1.16

WORKDIR /app

ADD data /app/data
COPY --from=builder /app/toolset /app/
COPY --from=builder /app/garble /go/bin/
CMD ["/app/toolset"]
