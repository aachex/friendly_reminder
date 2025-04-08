FROM golang:1.24.2-alpine3.21

WORKDIR /app
COPY . /app/

RUN apk update

# dependencies
RUN go get -d -u github.com/joho/godotenv
RUN go get -d -u github.com/golang-jwt/jwt/v5
RUN go get -d -u github.com/lib/pq

EXPOSE 8080

STOPSIGNAL SIGINT

CMD [ "go", "run", "/app/main.go" ]