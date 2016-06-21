# Command

## One-off request

### Read

#### Restful
- Verb: **GET**
- URI: /mb/tcp/**{fc}**/read
- query string: ?ip=**{ip}**&port=**{port}**&id=**{id}**&s=**{reg}**[&l=**{len}**]

|param|desc|type|range|example|optional|
|:--|:--|:--|:--|:--|:--|
|fc|read function code|int|[1,4]|1|-|
|ip|ip address|string|-| 127.0.0.1|-|  
|port|port number|int|[1,65535]|502|default 502|
|id|device id|int|[1, 253]|1|-|
|reg|register start addr|int|-|23|-|
|len|register length|int|-|20|fc1, fc2|

- params:
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
- URI: /mb/tcp/{fc}/write
- Example
- Response
- Example




## Polling request

### Add request
- HTTP verb: POST
- URI: /mb/tcp/{fc}/poll

## Filter request