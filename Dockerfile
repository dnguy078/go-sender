FROM golang:1.12
WORKDIR /appdir

# Fetch dependencies
COPY go.mod .
RUN go mod download

# # Build
COPY . .
RUN CGO_ENABLED=0 go build
RUN pwd && ls -lah

# todo move to alpine image final image
EXPOSE 4001 81

CMD ["./appdir/go-sender"]