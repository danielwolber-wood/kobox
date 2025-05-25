FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go mod tidy
RUN go mod download
RUN go build -o main ./cmd/Kobox-Mono

FROM alpine:latest
RUN apk --no-cache add ca-certificates pandoc
WORKDIR /root/
COPY --from=builder /app/main ./main
RUN mkdir -p /app/certs
EXPOSE 12332
CMD ["./main"]