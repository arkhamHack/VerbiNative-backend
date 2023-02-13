FROM golang as builder
RUN apk --no-cache add ca-certificates git
WORKDIR /build/api
COPY go.mod ./
RUN go mod download
COPY . ./RUN CGO_ENABLED=0 go build -o api
# post build stage
FROM alpine
WORKDIR /rootCOPY --from=builder /build/api/api .
EXPOSE 8080
CMD ["./api"]