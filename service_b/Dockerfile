# Usa a imagem base oficial do Golang 1.22
FROM golang:1.22

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos do módulo e o código para o contêiner
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Compila o aplicativo
RUN go build -o service_b .

# Exponha a porta que o serviço vai ouvir
EXPOSE 8081

# Comando para executar o serviço
CMD ["./service_b"]
