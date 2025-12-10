FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/main.go --parseDependency

RUN go build -o argent_api ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/argent_api .

COPY --from=builder /app/docs ./docs

EXPOSE 9000

CMD ["./argent_api"]
