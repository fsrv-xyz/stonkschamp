FROM golang:alpine@sha256:32c85006b1edf29c097514e0c81a33334aa1450685a885c10657ec756dbb7703 AS builder
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/app
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=builder /build/app /app/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/app"]

