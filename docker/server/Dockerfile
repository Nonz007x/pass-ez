FROM golang:1.22.2

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY . .

EXPOSE 3000

RUN go mod tidy