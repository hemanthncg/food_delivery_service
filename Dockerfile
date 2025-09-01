# syntax=docker/dockerfile:1
FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o app .
RUN chmod +x app

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
