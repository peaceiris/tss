FROM golang:1.18-alpine3.14 AS dev

ENV GOOS=linux
ENV GO111MODULE=on

RUN apk add --update-cache --no-cache bash && \
    go get github.com/magefile/mage
WORKDIR /go/src/github.com/peaceiris/tss

# ---

FROM dev AS builder

COPY . /go/src/github.com/peaceiris/tss/

RUN mage install

# ---

FROM alpine:3.17

COPY --from=builder /go/bin/tss /usr/bin/tss

CMD ["tss", "-h"]
