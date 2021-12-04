FROM golang:1.17-alpine as builder

WORKDIR /app

COPY cmd internal go.mod go.sum ./

RUN go build -o ./bin/secretsantabot ./cmd/secretsantabot/main.go
RUN go install github.com/golang-migrate/migrate/v4

FROM alpine:3.15
WORKDIR /app

RUN addgroup -S secretsantabot \
    && adduser -S secretsantabot -G secretsantabot -u 501
USER secretsantabot

COPY --from=builder /app/bin/secretsantabot ./secretsantabot
COPY --from=builder $GOPATH/bin/migrate ./migrate

CMD ["./secretsantabot"]
