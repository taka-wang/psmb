# Zero MQ Command Definition

# Table of contents

<!-- TOC depthFrom:2 depthTo:2 insertAnchor:false orderedList:false updateOnSave:true withLinks:true -->

- [0. Multipart message](#0-multipart-message)
- [1. One-off requests](#1-one-off-requests)
- [2. Polling requests](#2-polling-requests)
- [3. Filter requests](#3-filter-requests)

<!-- /TOC -->


## 0. Multipart message

We can compose a message out of several frames, and then receiver will receive all parts of a message, or none at all.
Thanks to the all-or-nothing characteristics, we can screen what we are interested from the first frame without parsing the whole JSON payload. 

>| Frame 1     |  Frame 2      |
>|:-----------:|:-------------:|
>| Method Name |  JSON Command |

---

## 1. One-off requests

**Data type**

>| type| description                            | args                                          | example                     | note    |
>|:----|:---------------------------------------|:----------------------------------------------|:----------------------------|:--------|
>| 1   | Raw register array (0xABCD hex array)  | -                                             | [0xABCD, 0x1234, 0xAB12]    | -       |
>| 2   | Hexadecimal string                     | -                                             | "112C004F12345678"          | -       |
>| 3   | Linearly scale uint16 to desired range | range: a (low), b (high), c (low), d (high)   | [22.34, 33.12, 44.56]       | -       |
>| 4   | uint16                                 | order: 1 (big-endian), 2 (little-endian)      | [123, 456, 789]             | -       |
>| 5   | int16                                  | order: 1 (big-endian), 2 (little-endian)      | [123, 456, 789]             | -       |
>| 6   | uint32                                 | order: 1 (ABCD), 2 (DCBA), 3 (BADC), 4 (CDAB) | [65538, 456, 789]           | len: 2x |
>| 7   | int32                                  | order: 1 (ABCD), 2 (DCBA), 3 (BADC), 4 (CDAB) | [65538, 456, 789]           | len: 2x |
>| 8   | float32                                | order: 1 (ABCD), 2 (DCBA), 3 (BADC), 4 (CDAB) | [22.34, 33.12, 44.56]       | len: 2x |

### 1.1 Read coil/register (**mbtcp.once.read**)
Command name: **mbtcp.once.read**

>| params   | description            | type          | range     | example           | required                                 |
>|:---------|:-----------------------|:--------------|:----------|:------------------|:-----------------------------------------|
>| from     | Service name           | string        | -         | "web"             | optional                                 |
>| tid      | Transaction ID         | integer       | int64     | 12345             | :heavy_check_mark:                       |
>| fc       | Function code          | integer       | [1,4]     | 1                 | :heavy_check_mark:                       |
>| ip       | IP address             | string        | -         | 127.0.0.1         | :heavy_check_mark:                       |
>| port     | Port number            | string        | [1,65535] | 502               | default: 502                             |
>| slave    | Slave id               | integer       | [1, 253]  | 1                 | :heavy_check_mark:                       |
>| addr     | Register start address | integer       | -         | 23                | :heavy_check_mark:                       |
>| len      | Bit/Register length    | integer       | -         | 20                | default: 1                               |
>| type     | Data type              | category      | [1,8]     | see below         | default: 1, **fc 3, 4 only**             |
>| order    | Endian                 | category      | [1,4]     | see below         | default: 1, **fc 3, 4 and type 4~8 only**|
>| range    | Scale range            | 4 floats      | -         | see below         | fc 3, 4 and type 3 only                  |
>| status   | Response status        | string        | -         | "ok"              | :heavy_check_mark:                       |
>| data     | Response value         | integer array |           | [1, 0, 24, 1]     | if success                               |
>| bytes    | Response byte array    | bytes array   | -         | [AB, 12, CD, ED]  | fc 3, 4 and type 2~8 only                |


#### 1.1.1 Services to PSMB

**Bits read (FC1, FC2)**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 1,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4
}
```

**Register read (FC3, FC4) - type 1, 2 (raw)**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 3,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "type": 1
}
```

**Register read (FC3, FC4) - type 3 (scale)**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 3,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "type": 3,
    "range": 
    {
        "a": 0,
        "b": 65535,
        "c": 100,
        "d": 500
    }
}
```

**Register read (FC3, FC4) - type 4, 5 (16-bit)**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 3,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "type": 4,
    "order": 1
}
```

**Register read (FC3, FC4) - type 6, 7, 8 (32-bit)**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 3,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "type": 6,
    "order": 3
}
```

#### 1.1.2 PSMB to Services

**Bits read (FC1, FC2)**

- Success:

```JavaScript
{
    "tid": 123456,
    "status": "ok",
    "data": [0,1,0,1,0,1]
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

**Register read (FC3, FC4) - type 1, 2 (raw)**

- Success - type 1 (RegisterArray):

```JavaScript
{
    "tid": 123456,
    "status": "ok",
    "type": 1,
    "bytes": [0XFF, 0X34, 0XAB],
    "data": [255, 1234, 789]
}
```

- Success - type 2 (HexString):

```JavaScript
{
    "tid": 123456,
    "status": "ok",
    "type": 2,
    "bytes": [0XFF, 0X34, 0XAB],
    "data": "112C004F12345678"
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

**Register read (FC3, FC4) - type 3 (scale)**

- Success - type 3 (Scale):

```JavaScript
{
    "tid": 123456,
    "status": "ok",
    "type": 3,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "data": [22.34, 33.12, 44.56]
}
```

- Fail - conversion fail:

```JavaScript
{
    "tid": 123456,
    "type": 3,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "status": "conversion fail"
}
```

- Fail - modbus fail:

```JavaScript
{
    "tid": 123456,
    "type": 3,
    "bytes": null,
    "status": "timeout"
}
```

**Register read (FC3, FC4) - type 4, 5 (16-bit)**

- Success - type 4, 5 (UInt16, Int16):

```JavaScript
{
    "tid": 123456,
    "status": "ok",
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "data": [255, 1234, 789]
}
```

- Fail - conversion fail:

```JavaScript
{
    "tid": 123456,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "status": "conversion fail"
}
```

- Fail - modbus fail:

```JavaScript
{
    "tid": 123456,
    "bytes": null,
    "status": "timeout"
}
```

**Register read (FC3, FC4) - type 6, 7, 8 (32-bit)**

- Success - type 6, 7 (UInt32, Int32):

```JavaScript
{
    "tid": 123456,
    "status": "ok",
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "data": [255, 1234, 789]
}
```

- Success - type 8 (Float32):

```JavaScript
{
    "tid": 123456,
    "status": "ok",
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "data": [22.34, 33.12, 44.56]
}
```

- Fail - Conversion fail:

```JavaScript
{
    "tid": 123456,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "status": "conversion fail"
}
```

- Fail - Modbus fail:

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

### 1.2 Write coil/register (**mbtcp.once.write**)
Command name: **mbtcp.once.write**

>| params   | description            | type          | range          | example        | required            |
>|:---------|:-----------------------|:--------------|:---------------|:---------------|:--------------------|
>| from     | service name           | string        | -              | "web"          | optional            |
>| tid      | transaction ID         | integer       | int64          | 12345          | :heavy_check_mark:  |
>| fc       | function code          | integer       | [1,4]          | 1              | :heavy_check_mark:  |
>| ip       | ip address             | string        | -              | 127.0.0.1      | :heavy_check_mark:  |
>| port     | port number            | string        | [1,65535]      | 502            | default: 502        |
>| slave    | slave id               | integer       | [1, 253]       | 1              | :heavy_check_mark:  |
>| addr     | register start address | integer       | -              | 23             | :heavy_check_mark:  |
>| len      | bit/register length    | integer       | -              | 20             | **FC15, 16 only**   |
>| **hex**  | hex/dec string flag    | bool          | [true, false]  | true           | **FC6, 16 only**    |
>| data(*)  | data to be write       | integer       | [0,1]          | 1              | **FC5 only**        |
>| data(**) | data to be write       | string        | hex/dec string | -              | **FC6, 16 only**    |
>| data(***)| data to be write       | integer array | bit array      | [1,1,0,1]      | **FC15 only**       |
>| status   | response status        | string        | -              | "ok"           | :heavy_check_mark:  |

#### 1.2.1 Services to PSMB

**bit write (FC5) - write single bit**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 5,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "data": 1
}
```

**register write (FC6) - write single register (dec)**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 6,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "hex": false,
    "data": "22"
}
```

**register write (FC6) - write single register (hex)**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 6,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "hex": true,
    "data": "ABCD"
}
```

**bits write (FC15) - write multiple bits**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 15,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "data": [1,0,1,0]
}
```

**registers write (FC16) - write multiple registers (dec)**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 16,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "hex": false,
    "data": "11,22,33,44"
}
```

**registers write (FC16) - write multiple registers (hex)**

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "fc" : 16,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "hex": true,
    "data": "ABCD1234EFAB1234"
}
```

#### 1.2.2 PSMB to Services

- Success:

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

### 1.3 Get TCP connection timeout (**mbtcp.timeout.read**)
Command name: **mbtcp.timeout.read**

>| params   | description            | type          | range          | example        | required                                     |
>|:---------|:-----------------------|:--------------|:---------------|:---------------|:---------------------------------------------|
>| from     | service name           | string        | -              | "web"          | optional                                     |
>| tid      | transaction ID         | integer       | int64          | 12345          | :heavy_check_mark:                           |
>| timeout  | timeout in usec        | integer       | [200000,~)	  | 210000         | if success                                   |
>| status   | response status        | string        | -              | "ok"           | :heavy_check_mark:                           |

#### 1.3.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 1468834562846161588
}
```

#### 1.3.2 PSMB to Services

- Success:

```JavaScript
{
    "tid": 1468834562846161588,
    "timeout": 212345,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 1468834562846161588,
    "status": "failed"
}
```

### 1.4 Set TCP connection timeout (**mbtcp.timeout.update**)
Command name: **mbtcp.timeout.update**

>| params   | description            | type          | range          | example        | required                                     |
>|:---------|:-----------------------|:--------------|:---------------|:---------------|:---------------------------------------------|
>| from     | service name           | string        | -              | "web"          | optional                                     |
>| tid      | transaction ID         | integer       | int64          | 12345          | :heavy_check_mark:                           |
>| timeout  | timeout in usec        | integer       | [200000,~)	  | 210000         | :heavy_check_mark:                           |
>| status   | response status        | string        | -              | "ok"           | :heavy_check_mark:                           |

#### 1.4.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 1468834562846161588,
    "timeout": 212345
}
```

#### 1.4.2 PSMB to Services

- Success:

```JavaScript
{
    "tid": 1468834562846161588,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 1468834562846161588,
    "status": "failed"
}
```

---

## 2. Polling requests


>| params       | description            | type          | range     | example           | required                                 |
>|:-------------|:-----------------------|:--------------|:----------|:------------------|:-----------------------------------------|
>| from         | Service name           | string        | -         | "web"             | optional                                 |
>| **name**     | poller name            | unique string | -         | "led_1"           | :heavy_check_mark:                       |
>| **interval** | polling interval in sec| integer       | [1~)      | 3                 | :heavy_check_mark:                       |
>|**enabled**   | polling enabled flag   | boolean       |true, false|true               |:heavy_check_mark:                        |
>| tid          | Transaction ID         | integer       | int64     | 12345             | :heavy_check_mark:                       |
>| fc           | Function code          | integer       | [1,4]     | 1                 | :heavy_check_mark:                       |
>| ip           | IP address             | string        | -         | 127.0.0.1         | :heavy_check_mark:                       |
>| port         | Port number            | string        | [1,65535] | 502               | default: 502                             |
>| slave        | Slave id               | integer       | [1, 253]  | 1                 | :heavy_check_mark:                       |
>| addr         | Register start address | integer       | -         | 23                | :heavy_check_mark:                       |
>| len          | Bit/Register length    | integer       | -         | 20                | default: 1                               |
>| type         | Data type              | category      | [1,8]     | see below         | default: 1, **fc 3, 4 only**             |
>| order        | Endian                 | category      | [1,4]     | see below         | default: 1, **fc 3, 4 and type 4~8 only**|
>| range        | Scale range            | 4 floats      | -         | see below         | fc 3, 4 and type 3 only                  |
>| status       | Response status        | string        | -         | "ok"              | :heavy_check_mark:                       |
>| data         | Response value         | integer array |           | [1, 0, 24, 1]     | if success                               |
>| bytes        | Response byte array    | bytes array   | -         | [AB, 12, CD, ED]  | fc 3, 4 and type 2~8 only                |



### 2.1 Add request (**mbtcp.poll.create**)
Command name: **mbtcp.poll.create**

#### 2.1.1 Services to PSMB

**Bits read (FC1, FC2)**

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "interval": 3,
    "enabled": true,
    "tid": 123456,
    "fc" : 1,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4
}
```

**Register read (FC3, FC4) - type 1, 2 (raw)**

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "interval": 3,
    "enabled": true,
    "tid": 123456,
    "fc" : 3,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "type": 1
}
```

**Register read (FC3, FC4) - type 3 (scale)**

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "interval": 3,
    "enabled": true,
    "tid": 123456,
    "fc" : 3,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "type": 3,
    "range": 
    {
        "a": 0,
        "b": 65535,
        "c": 100,
        "d": 500
    }
}
```

**Register read (FC3, FC4) - type 4, 5 (16-bit)**

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "interval": 3,
    "enabled": true,
    "tid": 123456,
    "fc" : 3,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "type": 4,
    "order": 1
}
```

**Register read (FC3, FC4) - type 6, 7, 8 (32-bit)**

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "interval": 3,
    "enabled": true,
    "tid": 123456,
    "fc" : 3,
    "ip": "192.168.0.1",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "type": 6,
    "order": 3
}
```

#### 2.1.2 PSMB to Services

**Bits read (FC1, FC2)**

- Success:

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

- Data:

```JavaScript
{
    "name": "led_1",
    "ts": 123456789,
    "status": "ok",
    "data": [0,1,0,1,0,1]
}
```

**Register read (FC3, FC4) - type 1, 2 (raw)**

- Success - type 1, 2:

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

- Data - type 1:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "ok",
    "type": 1,
    "data": [255, 1234, 789]
}
```

- Data - type 2:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "ok",
    "type": 2,
    "bytes": [0XFF, 0X34, 0XAB],
    "data": "112C004F12345678"
}
```

**Register read (FC3, FC4) - type 3 (scale)**

- Success - type 3:

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

- Data - type 3:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "ok",
    "type": 3,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "data": [22.34, 33.12, 44.56]
}
```

- Data - conversion fail - type 3:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "conversion fail",
    "type": 3,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34]
}
```

- Data - modbus fail - type 3:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "timeout",
    "type": 3,
    "bytes": null
}
```

**Register read (FC3, FC4) - type 4, 5 (16-bit)**

- Success - type 4, 5:

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

- Data - type 4,5:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "ok",
    "type": 4,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "data": [255, 1234, 789]
}
```

- Data - conversion fail - type 4, 5:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "conversion fail",
    "type": 4,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34]
}
```

- Data - modbus fail - type 4, 5:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "timeout",
    "type": 3,
    "bytes": null
}
```

**Register read (FC3, FC4) - type 6, 7, 8 (32-bit)**

- Success - type 6, 7, 8:

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

- Data - type 8:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "ok",
    "type": 8,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
    "data": [22.34, 33.12, 44.56]
}
```

- Data - conversion fail:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "conversion fail",
    "type": 6,
    "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34]
}
```

- Data - modbus fail:

```JavaScript
{
    "tid": 123456,
    "name": "led_1",
    "ts": 123456789,
    "status": "timeout",
    "type": 6,
    "bytes": null
}
```

### 2.2 Update request interval (**mbtcp.poll.update**)
Command name: **mbtcp.poll.update**

#### 2.2.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "interval": 3,
    "tid": 123456
}
```

#### 2.2.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

### 2.3 Read request status (**mbtcp.poll.read**)
Command name: **mbtcp.poll.read**

#### 2.3.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "tid": 123456
}
```

#### 2.3.2 PSMB to Services

- Success:

```JavaScript
{
    "name": "led_1",
    "tid": 123456,
    "fc": 1,
    "ip": "192.168.3.2",
    "port": "502",
    "slave": 22,
    "addr": 250,
    "len": 10,
    "interval" : 3,
    "status": "ok",
    "enabled": true,
    "type": xx,
    "order": yy,
    "range": {}
}
```

- Fail:

```JavaScript
{
    "name": "led_1",
    "tid": 123456,
    "status": "not exist"
}
```

### 2.4 Delete request (**mbtcp.poll.delete**)
Command name: **mbtcp.poll.delete**

#### 2.4.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "tid": 123456
}
```

#### 2.4.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "not exist"
}
```

### 2.5 Enable/Disable request (**mbtcp.poll.toggle**)
Command name: **mbtcp.poll.toggle**

#### 2.5.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "tid": 123456,
    "enabled": true
}
```

#### 2.5.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "not exist"
}
```

### 2.6 Read all requests status (**mbtcp.polls.read**)
Command name: **mbtcp.polls.read**

#### 2.6.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 123456,
}
```

#### 2.6.2 PSMB to Services

- success:

    ```JavaScript
    {
        "tid": 123456,
        "status": "ok",
        "polls": [
            {
                "name": "led_1",
                "fc": 1,
                "ip": "192.168.3.2",
                "port": "502",
                "slave": 22,
                "addr": 250,
                "len": 10,
                "interval" : 3,
                "status": "ok",
                "enabled": true
            },
            {
                "name": "led_2",
                "fc": 1,
                "ip": "192.168.3.2",
                "port": "502",
                "slave": 22,
                "addr": 250,
                "len": 10,
                "interval" : 3,
                "status": "ok",
                "enabled": true
            }]
    }
    ```

- fail:

    ```JavaScript
    {
        "tid": 123456,
        "status": "timeout"
    }
    ```

### 2.7 Delete all requests (**mbtcp.polls.delete**)
Command name: **mbtcp.polls.delete**

#### 2.7.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 123456
}
```

#### 2.7.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "not exist"
}
```

### 2.8 Enable/Disable all requests (**mbtcp.polls.toggle**)
Command name: **mbtcp.polls.toggle**

#### 2.8.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "enabled": true
}
```

#### 2.8.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "not exist"
}
```

### 2.9 Import requests (**mbtcp.polls.import**)
Command name: **mbtcp.polls.import**

#### 2.9.1 Services to PSMB
**TODO**

#### 2.9.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "not exist"
}
```

### 2.10 Export requests (**mbtcp.polls.export**)
Command name: **mbtcp.polls.export**

#### 2.10.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 123456
}
```

#### 2.10.2 PSMB to Services
**TODO**

### 2.11 Read history (**mbtcp.poll.history**)
Command name: **mbtcp.poll.history**

#### 2.11.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "tid": 123456
}
```

#### 2.11.2 PSMB to Services
**TODO**

---

## 3. Filter requests

**Equality Type**

>| type| description                  |  symbol | 
>|:----|:-----------------------------|:--------|
>| 0   | change                       |         |
>| 1   | greater than or equal (ge)   | >=      | 
>| 2   | greater than (gt)            | >       |
>| 3   | equal (eq)                   | ==      |
>| 4   | less than (lt)               | <       |
>| 5   | less than or equal (le)      | <=      |
>| 6   | inside range                 |         |
>| 7   | inside range (include)       |         |
>| 8   | outside range                |         |
>| 9   | outside range (include)      |         |


### 3.1 Add filter (**mbtcp.filter.create**)
Command name: **mbtcp.filter.create**

>| params       | description            | type          | range     | example     | required            |
>|:-------------|:-----------------------|:--------------|:----------|:------------|:--------------------|
>| from         | Service name           | string        | -         | "web"       | optional            |
>| **poll**     | poller name            | unique string | -         | "led_1"     | :heavy_check_mark:  |
>| **name**     | filter name            | unique string | -         | "filter_1"  | :heavy_check_mark:  |
>| tid          | Transaction ID         | integer       | int64     | 12345       | :heavy_check_mark:  |
>| type         | Comparison type        | category      | [0, 9]    |             | :heavy_check_mark:  |
>| arg          | value                  | array         |           |             | optional            |
>|**enabled**   | polling enabled flag   | boolean       |true, false|true         | :heavy_check_mark:  |
>| status       | Response status        | string        | -         | "ok"        | :heavy_check_mark:  |


#### 3.1.1 Services to PSMB

- type 0:

```JavaScript
{
    "from": "web",
    "poll": "led_1",
    "name": "filter_1",
    "tid": 123456,
    "enabled": true,
    "type": 0
}
```

- type 1, 2, 3, 4, 5:

```JavaScript
{
    "from": "web",
    "poll": "led_1",
    "name": "filter",
    "tid": 123456,
    "enabled": true,
    "type": 1,
    "arg": [1.0]
}
```

- type 6, 7, 8, 9:

```JavaScript
{
    "from": "web",
    "poll": "led_1",
    "name": "filter",
    "tid": 123456,
    "enabled": true,
    "type": 6,
    "arg": [1.0, 2.0]
}
```

#### 3.1.2 PSMB to Services

- Success:

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "fail"
}
```

### 3.2 Update filter (**mbtcp.filter.update**)
Command name: **mbtcp.filter.update**

#### 3.2.1 Services to PSMB

- type 0:

```JavaScript
{
    "from": "web",
    "poll": "led_1",
    "name": "filter_1",
    "tid": 123456,
    "enabled": true,
    "type": 0
}
```

- type 1, 2, 3, 4, 5:

```JavaScript
{
    "from": "web",
    "poll": "led_1",
    "name": "filter",
    "tid": 123456,
    "enabled": true,
    "type": 1,
    "arg": [1.0]
}
```

- type 6, 7, 8, 9:

```JavaScript
{
    "from": "web",
    "poll": "led_1",
    "name": "filter",
    "tid": 123456,
    "enabled": true,
    "type": 6,
    "arg": [1.0, 2.0]
}
```

#### 3.2.2 PSMB to Services

- Success:

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "tid": 123456,
    "status": "fail"
}
```

### 3.3 Read filter status (**mbtcp.filter.read**)
Command name: **mbtcp.filter.read**

#### 3.3.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "name": "led_1",
    "tid": 123456
}
```

#### 3.3.2 PSMB to Services

- success - type 0:

```JavaScript
{
    "from": "web",
    "poll": "led_1",
    "name": "filter_1",
    "tid": 123456,
    "enabled": true,
    "type": 0,
    "status": "ok"
}
```

- Success - type 1, 2, 3, 4, 5:

```JavaScript
{
    "from": "web",
    "poll": "led_1",
    "name": "filter",
    "tid": 123456,
    "enabled": true,
    "type": 1,
    "arg": [1.0],
    "status": "ok"
}
```

- Success - type 6, 7, 8, 9:

```JavaScript
{
    "from": "web",
    "poll": "led_1",
    "name": "filter",
    "tid": 123456,
    "enabled": true,
    "type": 6,
    "arg": [1.0, 2.0],
    "status": "ok"
}
```

- Fail:

```JavaScript
{
    "name": "filter_1",
    "tid": 123456,
    "status": "not exist"
}
```

### 3.4 Delete filter (**mbtcp.filter.delete**)
Command name: **mbtcp.filter.delete**

#### 3.4.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "name": "filter_1",
    "tid": 123456
}
```

#### 3.4.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "not exist"
}
```

### 3.5 Enable/Disable filter (**mbtcp.filter.toggle**)
Command name: **mbtcp.filter.toggle**

#### 3.5.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "name": "f_1",
    "tid": 123456,
    "enabled": true
}
```

#### 3.5.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "not exist"
}
```

### 3.6 Read all filters (**mbtcp.filters.read**)
Command name: **mbtcp.filters.read**

#### 3.6.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 123456
}
```

#### 3.6.2 PSMB to Services
**TODO**

### 3.7 Delete all filters (**mbtcp.filters.delete**)
Command name: **mbtcp.filters.delete**

#### 3.7.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 123456
}
```

#### 3.7.2 PSMB to Services
**TODO**

### 3.8 Enable/Disable all filters (**mbtcp.filters.toggle**)
Command name: **mbtcp.filters.toggle**

#### 3.8.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "enabled": true
}
```

#### 3.8.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

### 3.9 Import filters (**mbtcp.filters.import**)
Command name: **mbtcp.filters.import**

#### 3.9.1 Services to PSMB
**TODO**

#### 3.9.2 PSMB to Services

```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

### 3.10 Export filters (**mbtcp.filters.export**)
Command name: **mbtcp.filters.export**

#### 3.10.1 Services to PSMB

```JavaScript
{
    "from": "web",
    "tid": 123456
}
```

#### 3.10.2 PSMB to Services
**TODO**