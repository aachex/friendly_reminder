FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . /app/

RUN go build -o app.exe

EXPOSE 8080

STOPSIGNAL SIGINT

CMD [ "/app/app.exe" ]