FROM golang:1.20.4

WORKDIR /app

LABEL maintainer="ryan@hexa.mozmail.com"

COPY . .

RUN go build -o /app/go-playlist /app/pkg/main.go

ENTRYPOINT [ "/app/go-playlist", "serve" ]
