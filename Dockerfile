FROM golang:1.23.1-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /receipt-processor

EXPOSE 8080

CMD ["/receipt-processor"]
