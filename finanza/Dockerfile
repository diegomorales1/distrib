FROM golang:1.22.6

WORKDIR /Lab1

COPY . .

RUN go mod tidy
RUN go build -o finanza .

CMD ["./finanza"]