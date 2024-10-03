FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/app

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .
COPY internal/database/migrations /migrations

RUN apk add --no-cache postgresql-client

ENV DB_HOST=localhost
ENV DB_USER=postgres
ENV DB_PASSWORD=90814263
ENV DB_NAME=userBalanceAvito
ENV DB_PORT=5432

CMD ["sh", "-c", "psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f /migrations/1_create_user_account_managment.up.sql && ./main"]

EXPOSE 8000