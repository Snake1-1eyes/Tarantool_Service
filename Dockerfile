FROM golang:1.23-alpine3.21

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO=ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]