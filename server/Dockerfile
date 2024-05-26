FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY . .
COPY ../.env .

EXPOSE 3000

RUN go mod tidy

CMD ["air"]
