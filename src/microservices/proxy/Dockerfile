FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o proxy .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/proxy .
EXPOSE 8000
CMD ["./proxy"]
