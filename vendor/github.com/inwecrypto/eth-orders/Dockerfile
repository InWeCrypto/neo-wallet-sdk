FROM golang:1.9.2-stretch

LABEL maintainer="support@inwecrypto.com"

COPY . /go/src/github.com/inwecrypto/eth-orders

RUN go install github.com/inwecrypto/eth-orders/cmd/eth-orders && rm -rf /go/src

VOLUME ["/etc/inwecrypto/order/eth"]

WORKDIR /etc/inwecrypto/order/eth

EXPOSE 8000

CMD ["/go/bin/eth-orders","--conf","/etc/inwecrypto/order/eth/orders.json"]