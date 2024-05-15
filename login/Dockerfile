FROM golang:1.20.2

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", ".air.toml"]