FROM golang:1.20.5-alpine3.18

WORKDIR /app

ENV TZ="America/New_York"
LABEL maintainer="ryan@hexa.mozmail.com"

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-playlist /app/pkg/main.go

ENTRYPOINT [ "/app/go-playlist", "serve" ]
