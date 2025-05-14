FROM golang:1.23-alpine AS builder

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/chains-auth

FROM alpine:latest

WORKDIR /service

COPY --from=builder /service/main .
COPY config_docker.yaml ./config_docker.yaml

ENV KV_VIPER_FILE=/service/config_docker.yaml

EXPOSE 8001

CMD ["./main", "run", "service"]
