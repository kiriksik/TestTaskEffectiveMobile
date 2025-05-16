FROM golang:alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]