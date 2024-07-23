FROM golang:1.22.5-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o service ./cmd/lru-cache

FROM alpine:latest

COPY --from=build /app/service /app/service

EXPOSE 8080

ENTRYPOINT ["/app/service"]
