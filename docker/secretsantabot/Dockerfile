FROM golang:1.17-alpine as builder

WORKDIR /app

COPY src .
RUN go build -o ./bin/secretsantabot ./cmd/secretsantabot/main.go

FROM alpine:3.15
WORKDIR /app

RUN addgroup -S secretsantabot \
    && adduser -S secretsantabot -G secretsantabot -u 501
USER secretsantabot

COPY --from=builder /app/bin/secretsantabot ./secretsantabot

CMD ["./secretsantabot"]
