# Dockerfile para servicoA
FROM golang:1.16

WORKDIR /app

# Copia os arquivos go.mod e go.sum
COPY go.mod ./
COPY go.sum ./

# Instala as dependências
RUN go mod download

# Copia o restante dos arquivos
COPY . .

# Compila o serviço
RUN go build -o servicoA servicoA.go

# Comando de execução
CMD ["./servicoA"]
