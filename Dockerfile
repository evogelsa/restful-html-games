# syntax=docker/dockerfile:1

FROM golang:1.14-alpine

WORKDIR /app

COPY src/ ./

RUN go mod download

RUN go build -o /restful-html-games

CMD [ "/restful-html-games" ]
