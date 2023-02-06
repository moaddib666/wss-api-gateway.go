# WSS API Gateway

This service transform websocket messages to pub/sub communication via RMQ.
## Architecture

### Logical view
[<img src=".docs/WSSApiGW.png">](http://google.com.au/)

### Physical view

[<img src=".docs/WSSApiGW-phisical.png">](http://google.com.au/)

## Local stand
1) Run RMQ bitnami image used in example: `docker compose up -d`
2) Run application `go run main.go`
3) Login to RMQ dashboard via `http://127.0.0.1:15672/`
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
- Payload:
 ```json
  {
  "sub": "1234567890",
  "user": "John Snow",
  "iat": 1516239022
  }
  ```

**Limitation**: for now only one connection allowed for 1 user.

### Auth 
   Auth module is under definition.
   But main concept is to perform auth request to the Identity Provider and get cleintId from the token

### Outbox
In RMQ transport represented as `topic` exchange by default name is `ApiGatewayOutbox`

#### RMQ Messages
On Client Connected
- Exchange:	`ApiGatewayOutbox`
- Routing Key: ``
- Properties
  - headers:
    - `recipient`: `*`
    - `sender`:	`ApiGateway`
  - Payload
     ```json 
     {
        "metadata":{
           "entity_id":"TestClient",
           "entity_name":"client",
           "publisher_name":"ApiGateway",
           "event_name":"connected",
           "created":"2023-02-06T20:31:06+02:00"
        },
        "payload":{
         
        }
     }
     ```

On Client Disconnected
- Exchange:	`ApiGatewayOutbox`
- Routing Key: ``
- Properties
  - headers:
    - `recipient`: `*`
    - `sender`:	`ApiGateway`
  - Payload
     ```json 
     {
        "metadata":{
           "entity_id":"TestClient",
           "entity_name":"client",
           "publisher_name":"ApiGateway",
           "event_name":"disconnected",
           "created":"2023-02-06T20:31:06+02:00"
        },
        "payload":{
         
        }
     }
     ```

### Inbox

In RMQ transport represented as `direct` exchange by default name is `ApiGatewayInbox`
Require that message headers to be set like:
- headers:
  - `recipient`: `ClientId`
  - `sender`:	`MyAwesomeMss`
- Payload:
  - Accept serialized aka `json`,`yaml` and plain text messages.






