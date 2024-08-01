FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY main.go ./
RUN go build -v -o ./app

EXPOSE 8080

CMD ["./app"]
