FROM golang:1.19-alpine

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

CMD ["/app/main"]