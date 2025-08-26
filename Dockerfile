# Base stage
FROM golang:1.24.5-alpine AS base
WORKDIR /app

RUN apk add --no-cache make git
# install Sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sqlc generate
# Dev stage
FROM base AS dev
RUN go install github.com/air-verse/air@latest
EXPOSE ${PORT}
CMD ["air", "-c", ".air.toml"]

# Prod stage
FROM base AS prod
RUN go build -o /app/main ./cmd/api/main.go
EXPOSE ${PORT}
CMD ["/app/main"]
