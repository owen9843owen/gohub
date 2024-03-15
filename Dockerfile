FROM golang:1.21

WORKDIR /server
COPY . /server

RUN go mod tidy
RUN go build -v -o app

EXPOSE 8088

ENTRYPOINT ./app