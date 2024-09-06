FROM golang:alpine@sha256:c76fc1ddd43f2f237b80a55e5f10f97ab2aea53f9a44c74d2495ac9536e8fcb4 AS builder
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/app
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=builder /build/app /app/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/app"]

