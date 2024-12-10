# Dockerfile for Servico B
FROM golang:1.16

WORKDIR /app

COPY . .

RUN go build -o servicoB servicoB.go

CMD ["./servicoB"]
