# syntax=docker/dockerfile:1

FROM golang:1.21.5

WORKDIR /activitypub-playground

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /ap-server

EXPOSE 8080

# Run
CMD ["/ap-server"]
