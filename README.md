# goLangBotTelegram

# start in background:
docker-compose up --build -d 

# down:
docker-compose down

# Connection in browser
http://localhost:8080

# Telegram token 

create .env file with key toke_telegram and value bot token
example: toke_telegram=bottoken

# Hello
    Port: 8080

# Produtos
    Port: 7070

# Suporte
    Port: 6060

# Criando novo package
    go mod init  <package-name>
    go mod tidy

# Build image dockerfile por modulo
    docker build -f Dockerfile.suporte -t my-suporte-image .
