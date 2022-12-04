FROM golang:alpine
WORKDIR /app
COPY . /app
RUN go mod tidy
RUN go build main.go
CMD ["./main"]