FROM golang:1.22

WORKDIR /app

COPY . .
RUN go mod download && go mod verify

RUN go build -v -o ./app/main.go

EXPOSE 8080

CMD ["./app/main.go"]
