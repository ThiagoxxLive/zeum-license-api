FROM golang:1.25-alpine

RUN apk add --no-cache git

WORKDIR /app

EXPOSE 8080

CMD ["sh", "-c", "go mod tidy && go run ./cmd/api"]
