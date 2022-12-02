FROM golang:1.16 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN apt-get update -y \
  && apt-get install --no-install-recommends -y librocksdb-dev \
  && apt-get clean

RUN go mod tidy && CGO_ENABLED=1 go build -a -o etcd-carry main.go

FROM alpine:3.13.5

WORKDIR /
COPY --from=builder /workspace/etcd-carry .

ENTRYPOINT ["/etcd-carry"]
