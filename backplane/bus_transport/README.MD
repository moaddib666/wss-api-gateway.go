# RMQ transport
Used as backplane for api gateway

### Inbox Queue
Durable queue witch used as message broker.
Microservices that want to communicate through WSS need to publish their massages to `inbox` queue

Default name is `ApiGatewayInbox`

### Outbox Exchange
Exchange that accepts messages from api gateway.
Microservices that wants to listen events from WSS need to create separate `service` queue and bind it the `outbox` exchange

Default name is: `ApiGatewayOutbox` 

### Message Protocol

#### Headers
- `recipient` - client id that should receive message, in case header empty or omitted - message will be skipped
- `sender`    - application id from message has been sent, in case header empty or omitted - `unknown` will be used as default value

#### Body
Could be text/plain or text/json or other serialized message
Example:
```json
{"foo": "bar"}
```
