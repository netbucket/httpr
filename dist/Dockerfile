# Run a clean build
FROM golang:1.11 as builder

RUN go get github.com/spf13/cobra
RUN go get github.com/netbucket/httpr

WORKDIR /go/src/github.com/netbucket/httpr
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o httpr

# Reduce the image size and add TLS certs
FROM alpine:3.7
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /
COPY --from=builder /go/src/github.com/netbucket/httpr/httpr .

EXPOSE 80

ENTRYPOINT ["/httpr"]
CMD ["log", "-e", "-p", "-s=:80"]
