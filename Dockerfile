FROM golang:1.24.5-alpine AS build

WORKDIR /app

RUN apk add --no-cache make git
# install Sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sqlc generate

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/web ./web
COPY --from=build /app/.env ./.env
EXPOSE ${PORT}
CMD ["./main"]
