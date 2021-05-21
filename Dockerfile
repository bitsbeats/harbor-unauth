FROM golang:1.16-alpine AS builder

ADD . /src
WORKDIR /src
RUN go build -o unauth .

# ---

FROM alpine:3.13

RUN apk add --no-cache ca-certificates
COPY --from=builder /src/unauth /usr/local/bin/unauth
ENTRYPOINT ["/usr/local/bin/unauth"]
