FROM golang:alpine3.19

USER root
WORKDIR /app

EXPOSE 8080

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build cmd/api/main.go

ENTRYPOINT [ "./main" ]