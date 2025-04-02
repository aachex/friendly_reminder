FROM golang:1.24.0-alpine3.21

WORKDIR /app
COPY . /app/

RUN apk update && apk add git build-base sqlite
RUN go get -d -u github.com/mattn/go-sqlite3
RUN go get -d -u github.com/joho/godotenv
RUN go get -d -u github.com/golang-jwt/jwt/v5

EXPOSE 8080

STOPSIGNAL SIGINT

CMD [ "go", "run", "/app/main.go" ]