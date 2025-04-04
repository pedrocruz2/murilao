FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN mkdir -p /app/uploads

# Inicializar módulo Go e adicionar dependências
RUN rm -f go.mod go.sum && \
    go mod init app && \
    go get github.com/gin-gonic/gin@v1.8.1 && \
    go get github.com/gin-contrib/cors@v1.4.0 && \
    go get github.com/golang-jwt/jwt/v5@v5.0.0 && \
    go get github.com/google/uuid@v1.3.0 && \
    go get github.com/joho/godotenv@v1.5.1 && \
    go get golang.org/x/crypto@v0.9.0 && \
    go get gorm.io/driver/postgres@v1.5.0 && \
    go get gorm.io/gorm@v1.25.0 && \
    go mod tidy

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]