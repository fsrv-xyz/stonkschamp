FROM golang:alpine@sha256:b036c52b3bcc8e4e31be19a7a902bb9897b2bf18028f40fd306a9778bab5771c AS builder
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/app
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=builder /build/app /app/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/app"]

