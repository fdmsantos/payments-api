FROM golang:latest
MAINTAINER Fabio Santos <fabiodanielmonteirosantos@gmail.com>

RUN mkdir api-payments
WORKDIR api-payments
COPY . .

RUN go build -o build/payments-api .

EXPOSE 8000

ENTRYPOINT ["./build/payments-api"]