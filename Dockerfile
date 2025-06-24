FROM golang:1.24-alpine

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .


RUN go build -v -o /usr/local/bin/app ./cmd/main.go

EXPOSE 8000

CMD ["app"]