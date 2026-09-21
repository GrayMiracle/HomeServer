FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o homeserver .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
RUN mkdir -p data
COPY --from=builder /app/homeserver .
COPY config.json .
EXPOSE 3164
CMD ["./homeserver"]