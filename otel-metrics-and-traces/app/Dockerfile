FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bootstrap

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/bootstrap /bootstrap
CMD ["/bootstrap"]