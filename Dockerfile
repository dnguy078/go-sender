FROM golang:1.12

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# todo: copy binary into a scratch image instead
EXPOSE 4001

CMD ["./go-sender"]