# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/

# Run Stage
FROM golang:1.26.1-bookworm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps

COPY --from=builder /app/main .

CMD ["./main"]