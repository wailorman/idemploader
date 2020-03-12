# Start from the latest golang base image
FROM golang:1.13

LABEL maintainer="Sergey Popov <sergey.popov.w@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
ADD sample.env .env

RUN go build -o main cmd/server/main.go

EXPOSE 8080

ENV IDEMPLOADER_PORT 3000

CMD ["./main"]
