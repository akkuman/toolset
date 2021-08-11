FROM golang:1.16 as builder

WORKDIR /src

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct

COPY ./go.mod ./go.sum /src/
RUN go get -u github.com/swaggo/swag/cmd/swag && go mod download
ADD . /src/
RUN swag init && go mod tidy && go build -o /app/toolset .


FROM golang:1.16

WORKDIR /app

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    GO111MODULE=on go get mvdan.cc/garble

ADD data /app/data
COPY --from=builder /app/toolset /app/
CMD ["/app/toolset"]
