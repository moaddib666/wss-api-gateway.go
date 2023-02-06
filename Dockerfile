FROM golang as builder

WORKDIR /go/src/github.com/moaddib666/wss-api-gateway.go

COPY . .

RUN go get .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o apiGateway .

FROM scratch

ENV RMQ_DSN=amqp://user:bitnami@rabbitmq:5672/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /bin/

COPY --from=builder /go/src/github.com/moaddib666/wss-api-gateway.go/apiGateway .

CMD [ "./apiGateway" ]

EXPOSE 8080