FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o loadtester .

ENTRYPOINT ["./loadtester"]