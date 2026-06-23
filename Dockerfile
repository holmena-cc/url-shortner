FROM golang:1.24.5-alpine AS build

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main cmd/api/main.go

# ---

FROM alpine:3.20.1 AS prod

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=build /app/main  ./main
COPY --from=build /app/web   ./web

EXPOSE 8080

CMD ["./main"]