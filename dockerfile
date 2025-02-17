FROM golang:1.23.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o news-app .

FROM alpine:latest

WORKDIR /usr/local/bin

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/news-app .

RUN chmod +x news-app

COPY config.env /usr/local/bin/config.env

COPY migrations /usr/local/bin/migrations

CMD ["news-app"]
