# Dockerfile for slow_api
FROM golang:1.24-alpine

WORKDIR /app

COPY ./slow_api .

RUN go build -o slow main.go

EXPOSE 8081

CMD ["./slow"]
