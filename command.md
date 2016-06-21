# Command

## One-off request

### Read

#### RESTFUL
- Verb: GET
- URI: /mb/tcp/{fc}/get?ip={ip}&port={port}&id={id}&s={reg}[&l={len}]
    - {fc}: [1,2,3,4]: int
    - {ip}: ip address: string
    - {port}: port number: int, optional
    - {id}: device id: int
    - {reg}: register start address: int
    - {len}: register length: int, optional for fc1, fc2
- EXAMPLE
- RESPONSE
- EXAMPLE

#### Websocket
##### REQ
##### RES

#### MQTT

#### DDS

#### ZMQ

---

### Write

#### RESTFUL
- Verb: POST
- URI: /mb/tcp/{fc}/set
- Example
- Response
- Example




## Polling request

### Add request
- HTTP verb: POST
- URI: /mb/tcp/{fc}/poll

## Filter request