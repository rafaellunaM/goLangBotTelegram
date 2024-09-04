FROM golang:1.22

WORKDIR /app

COPY . .
RUN go mod download && go mod verify
RUN go get -u github.com/go-telegram/bot
RUN go get github.com/joho/godotenv

RUN go build -v -o ./app/main.go

EXPOSE 5050

ENTRYPOINT ["./app/main.go"]

