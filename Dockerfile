FROM golang:1-alpine as builder

RUN set -x \
    && apk add --no-cache \
        git \
        curl \
    && mkdir -p /tmp \
    && mkdir -p /fixtures \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/01.txt \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/02.txt \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/03.txt \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/04.txt \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/05.txt \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/06.txt \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/07.txt \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/08.txt \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/09.txt \
    && curl -s https://loripsum.net/api/1/verylong/plaintext | awk 'NF' - | cat > /fixtures/10.txt

COPY . /tzhash

WORKDIR /tzhash

# https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away
# go build -mod=vendor
RUN set -x \
    && export CGO_ENABLED=0 \
    && go build -mod=vendor -o /go/bin/homo ./cmd/homo/main.go

# Executable image
FROM alpine:3.11

WORKDIR /fixtures

COPY --from=builder /fixtures    /fixtures
COPY --from=builder /go/bin/homo /usr/local/sbin/homo
