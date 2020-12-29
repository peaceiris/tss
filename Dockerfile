FROM golang:1.15-alpine3.12 AS dev

ENV GOOS=linux
ENV GO111MODULE=on

RUN apk add --no-cache bash && \
    go get github.com/magefile/mage
WORKDIR /go/src/github.com/peaceiris/tss

# ---

FROM dev AS builder

COPY . /go/src/github.com/peaceiris/tss/

RUN mage install

# ---

FROM alpine:3.12

COPY --from=builder /go/bin/tss /usr/bin/tss

CMD ["tss", "-h"]
