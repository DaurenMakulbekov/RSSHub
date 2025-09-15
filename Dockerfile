FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o rsshub ./cmd

FROM alpine

COPY --from=builder /app/rsshub /rsshub

CMD ["./rsshub"]
