FROM golang:1.17-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o ./bin/secretsantabot ./cmd/secretsantabot/main.go
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.1

FROM alpine:3.15
WORKDIR /app

RUN addgroup -S secretsantabot \
    && adduser -S secretsantabot -G secretsantabot -u 501
USER secretsantabot

COPY --from=builder /app/bin/secretsantabot ./secretsantabot
COPY --from=builder /go/bin/migrate ./migrate
COPY --from=builder /app/migrations ./migrations

CMD ["./secretsantabot"]
