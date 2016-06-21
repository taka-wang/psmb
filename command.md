# Command

## One-off request

### Read

#### RESTFUL
- Verb: GET
- URI: /mb/tcp/**{fc}**/get?ip=**{ip}**&port=**{port}**&id=**{id}**&s=**{reg}**[&l=**{len}**]

|paras|desc|type|range|example|optional|
|:--|:--|:--|:--|:--|:--|
|fc|function code|int|[1,4]|1||
|ip|ip address|string|| 127.0.0.1||  
|port|port number|int|[1,65535]|502|yes|
|id|device id|int|[1, 253]|1||
|reg|register start addr|int||23||
|len|register length|int||20|yes|

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