FROM golang:1.18 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/

#ENV GOPROXY="https://goproxy.cn"

RUN go mod tidy && CGO_ENABLED=1 go build -a -o etcd-carry main.go

FROM ubuntu:20.04

WORKDIR /
COPY --from=builder /workspace/etcd-carry .

CMD ["/etcd-carry"]
