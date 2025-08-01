FROM golang:1.22-alpine AS builder


WORKDIR /app


RUN apk add --no-cache git


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o api ./cmd/api


FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/api .


EXPOSE 8080


CMD ["./api"]
