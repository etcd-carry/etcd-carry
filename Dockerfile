FROM ubuntu:20.04 as builder

CMD ["bash"]

RUN set -eux \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
       ca-certificates netbase wget git mercurial openssh-client subversion procps g++ gcc libc6-dev make pkg-config librocksdb-dev \
    && if ! command -v gpg > /dev/null; then \
         apt-get update; apt-get install -y --no-install-recommends gnupg dirmngr; \
       fi \
    && rm -rf /var/lib/apt/lists/* \
    && url='https://dl.google.com/go/go1.18.10.linux-amd64.tar.gz' \
    && wget -O go.tgz.asc "$url.asc" \
    && wget -O go.tgz "$url" --progress=dot:giga \
    && GNUPGHOME="$(mktemp -d)" \
    && export GNUPGHOME \
    && gpg --batch --keyserver keyserver.ubuntu.com --recv-keys 'EB4C 1BFD 4F04 2F6D DDCC  EC91 7721 F63B D38B 4796' \
    && gpg --batch --keyserver keyserver.ubuntu.com --recv-keys '2F52 8D36 D67B 69ED F998  D857 78BD 6547 3CB3 BD13' \
    && gpg --batch --verify go.tgz.asc go.tgz \
    && gpgconf --kill all \
    && rm -rf "$GNUPGHOME" go.tgz.asc \
    && tar -C /usr/local -xzf go.tgz \
    && rm go.tgz

ENV PATH=$PATH:/usr/local/go/bin
ENV GOLANG_VERSION=1.18.10
ENV GOPATH=/go
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN go mod tidy && CGO_ENABLED=1 go build -a -o etcd-carry main.go

FROM ubuntu:20.04

RUN apt-get update -y \
  && apt-get install --no-install-recommends -y librocksdb-dev \
  && apt-get clean && rm -rf /var/log/*log /var/lib/apt/lists/* /var/log/apt/* /var/lib/dpkg/*-old /var/cache/debconf/*-old

WORKDIR /
COPY --from=builder /workspace/etcd-carry .

CMD ["/etcd-carry"]
