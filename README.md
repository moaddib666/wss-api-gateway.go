# Margay - Websocket API Gateway

This service transform websocket messages to pub/sub communication via RMQ.

Our gateway - `Margay` like a small wild cat, consume minimum resources but provide flexible and scalable solution.


## Architecture

### Logical view
[<img src=".docs/WSSApiGW.png">](http://google.com.au/)

### Physical view

[<img src=".docs/WSSApiGW-phisical.png">](http://google.com.au/)

## Local stand
1) Run RMQ bitnami image used in example: `docker compose up -d`
2) Setup transport dsn via env `MARGAY_TRANSPORT_DSN`
3) Run application `go run main.go`
4) Login to RMQ dashboard via `http://127.0.0.1:15672/`
   - user: `user`
   - password: `bitnami`

## Protocol

### WSS connection

First request mast have `Autorization` header with JWT that represent user id

Example: 

`Autorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlciI6IkpvaG4gU25vdyIsImlhdCI6MTUxNjIzOTAyMn0.UrZ7NLiOUVVUerPuKouyr7NQz8XNb3jMrS_m0_07o8A`

Token Decoded:
- Header:
 ```json
  {
  "alg": "HS256",
  "typ": "JWT"
  }
  ```
- Claims: https://www.rfc-editor.org/rfc/rfc7519#section-4.1
 ```json
  {
  "aud": "localhost", 
  "exp": "1.675816912e+09",
  "iat": "1.675730512e+09",
  "iss": "IdentityServiceLocal",
  "jti": "Jhohn Snow",
  "sub": "client"
  }
  ```

**Limitation**: for now only one connection allowed for 1 user.

### Auth 
   Auth module is under definition.
   But main concept is to perform auth request to the Identity Provider and get cleintId from the token

### Outbox
In RMQ transport represented as `topic` exchange by default name is `MargayGatewayOutbox`

#### RMQ Messages
On Client Connected
- Exchange:	`MargayGatewayOutbox`
- Routing Key: ``
- Properties
  - headers:
    - `recipient`: `*`
    - `sender`:	`MargayGateway`
  - Payload
  ```json
  {
     "meta":{
        "object_id":"TestClient",
        "object_type":"client",
        "publisher":"MargayGateway",
        "event":"connected",
        "created":"2023-02-06T20:51:49+02:00"
     },
     "data":{
        
     }
  }
  ```

On Client Disconnected
- Exchange:	`MargayGatewayOutbox`
- Routing Key: ``
- Properties
  - headers:
    - `recipient`: `*`
    - `sender`:	`MargayGateway`
  - Payload
  ```json
  {
     "meta":{
        "object_id":"TestClient",
        "object_type":"client",
        "publisher":"MargayGateway",
        "event":"disconnected",
        "created":"2023-02-06T20:51:49+02:00"
     },
     "data":{
        
     }
  }
  ```
### Inbox

In RMQ transport represented as `direct` exchange by default name is `MargayGatewayInbox`
Require that message headers to be set like:
- headers:
  - `recipient`: `ClientId`
  - `sender`:	`MyAwesomeMss`
- Payload:
  - Accept serialized aka `json`,`yaml` and plain text messages.






