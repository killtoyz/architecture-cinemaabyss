FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY main.go ./
# go.sum в репозитории нет — tidy подтянет зависимости и создаст его при сборке.
RUN go mod tidy && CGO_ENABLED=0 go build -o events .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/events .
EXPOSE 8082
CMD ["./events"]
