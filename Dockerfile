FROM golang:1.12
WORKDIR /appdir

# Fetch dependencies
COPY go.mod .
RUN go mod download

# # Build
COPY . .
RUN go build

# move to alpine image
CMD ["./appdir/go-sender"]