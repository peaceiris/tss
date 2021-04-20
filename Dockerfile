FROM golang:1.16-alpine3.13 AS dev

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

FROM alpine:3.13

COPY --from=builder /go/bin/tss /usr/bin/tss

CMD ["tss", "-h"]
