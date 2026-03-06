# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN apk add --no-cache git
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates curl postgresql-client tzdata

WORKDIR /app

COPY --from=builder /server .
COPY db/migrations ./db/migrations
COPY db/seeds ./db/seeds

# Install migrate tool
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
