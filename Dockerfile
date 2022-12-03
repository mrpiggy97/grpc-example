FROM golang:alpine
WORKDIR /grcp-server
COPY . /grpc-server
RUN go mod tidy
RUN go build main.go
CMD ["./main"]