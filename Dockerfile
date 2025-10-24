# Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copia apenas go.mod/go.sum e baixa dependências
COPY go.mod go.sum ./
RUN go mod download

# Build só do código Go
COPY . .
RUN go build -o main .
RUN go build -o setup ./cmd/start/setup.go

# Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/setup .
# Copia a pasta docs inteira para o runtime
COPY docs ./docs

# Remove apenas o docs.go que não é necessário
RUN rm -f ./docs/docs.go

# Dá permissão de execução para os binários
RUN chmod +x ./main
RUN chmod +x ./setup

# Porta que a aplicação irá rodar
EXPOSE 8080

# Comando para rodar a aplicação
# CMD ["./main"]

# Comando para rodar a aplicação com dados de teste
CMD ["sh", "-c", "./setup && ./main"]