FROM golang:alpine@sha256:b4766d6a4baaa4c4cd88a3b9e0373d071e523fd8b42a3256f3f4c4a4fa887609 AS builder
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/app
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=builder /build/app /app/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/app"]

