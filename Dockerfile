FROM golang:1.24.0-alpine3.21

WORKDIR /app

COPY . /app/

RUN apk update && apk add git sqlite

CMD [ "go", "run", "/app/main.go" ]