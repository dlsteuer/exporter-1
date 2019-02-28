FROM golang:1.11-alpine as builder
WORKDIR /go/src/github.com/battlesnakeio/exporter/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go install -installsuffix cgo ./cmd/...

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY ./render/assets /app/render/assets
COPY --from=builder /go/bin/ /bin/
CMD ["/bin/exporter"]
