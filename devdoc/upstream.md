# Zero MQ Command Definition

# Table of contents

<!-- TOC depthFrom:2 depthTo:2 insertAnchor:false orderedList:false updateOnSave:true withLinks:true -->

- [0. Multipart message](#0-multipart-message)
- [1. One-off requests](#1-one-off-requests)
	- [1.1 Read coil/register](#11-read-coilregister)
		- [1.1.1 services to psmb](#111-services-to-psmb)
		- [1.1.2 psmb to services](#112-psmb-to-services)
	- [1.2 Write coil/register](#12-write-coilregister)
		- [1.2.1 services to psmb](#121-services-to-psmb)
		- [1.2.2 psmb to services](#122-psmb-to-services)
	- [1.3 Get TCP connection timeout](#13-get-tcp-connection-timeout)
		- [1.3.1 services to psmb](#131-services-to-psmb)
		- [1.3.2 psmb to services](#132-psmb-to-services)
	- [1.4 Set TCP connection timeout](#14-set-tcp-connection-timeout)
		- [1.4.1 services to psmb](#141-services-to-psmb)
		- [1.4.2 psmb to services](#142-psmb-to-services)
- [2. Polling requests](#2-polling-requests)
	- [2.1 Add request](#21-add-request)
	- [2.2 Update request (interval)](#22-update-request-interval)
	- [2.3 Read request status](#23-read-request-status)
	- [2.4 Delete request](#24-delete-request)
	- [2.5 Enable/Disable request](#25-enabledisable-request)
	- [2.6 Read all requests status](#26-read-all-requests-status)
	- [2.7 Delete all requests](#27-delete-all-requests)
	- [2.8 Enable/Disable all requests](#28-enabledisable-all-requests)
	- [2.9 Import requests](#29-import-requests)
	- [2.10 Export requests](#210-export-requests)
	- [2.11 Read history](#211-read-history)
- [3. Filter requests](#3-filter-requests)
	- [3.1 Add filter](#31-add-filter)
	- [3.2 Update filter](#32-update-filter)
	- [3.3 Read filter status](#33-read-filter-status)
	- [3.4 Delete filter](#34-delete-filter)
	- [3.5 Enable/Disable filter](#35-enabledisable-filter)
	- [3.6 Read all filters](#36-read-all-filters)
	- [3.7 Delete all filters](#37-delete-all-filters)
	- [3.8 Enable/Disable all filters](#38-enabledisable-all-filters)
	- [3.9 Import filters](#39-import-filters)
	- [3.10 Export filters](#310-export-filters)

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
>| 1   | raw register array (0xABCD hex array)  | -                                             | [0xABCD, 0x1234, 0xAB12]    | -       |
>| 2   | hexadecimal string                     | -                                             | "112C004F12345678"          | -       |
>| 3   | linearly scale uint16 to desired range | range: a (low), b (high), c (low), d (high)   | [22.34, 33.12, 44.56]       | -       |
>| 4   | uint16                                 | order: 1 (big-endian), 2 (little-endian)      | [123, 456, 789]             | -       |
>| 5   | int16                                  | order: 1 (big-endian), 2 (little-endian)      | [123, 456, 789]             | -       |
>| 6   | uint32                                 | order: 1 (ABCD), 2 (DCBA), 3 (BADC), 4 (CDAB) | [65538, 456, 789]           | len: 2x |
>| 7   | int32                                  | order: 1 (ABCD), 2 (DCBA), 3 (BADC), 4 (CDAB) | [65538, 456, 789]           | len: 2x |
>| 8   | float32                                | order: 1 (ABCD), 2 (DCBA), 3 (BADC), 4 (CDAB) | [22.34, 33.12, 44.56]       | len: 2x |

### 1.1 Read coil/register
Command name: **mbtcp.once.read**

>| params   | description            | type          | range     | example           | required                                 |
>|:---------|:-----------------------|:--------------|:----------|:------------------|:-----------------------------------------|
>| from     | service name           | string        | -         | "web"             | optional                                 |
>| tid      | transaction ID         | integer       | uint64    | 12345             | :heavy_check_mark:                       |
>| fc       | function code          | integer       | [1,4]     | 1                 | :heavy_check_mark:                       |
>| ip       | ip address             | string        | -         | 127.0.0.1         | :heavy_check_mark:                       |
>| port     | port number            | string        | [1,65535] | 502               | default: 502                             |
>| slave    | slave id               | integer       | [1, 253]  | 1                 | :heavy_check_mark:                       |
>| addr     | register start address | integer       | -         | 23                | :heavy_check_mark:                       |
>| len      | bit/register length    | integer       | -         | 20                | default: 1                               |
>| type     | data type              | category      | [1,8]     | see below         | default: 1, **fc 3, 4 only**             |
>| order    | endian                 | category      | [1,4]     | see below         | default: 1, **fc 3, 4 and type 4~8 only**|
>| range    | scale range            | 4 floats      | -         | see below         | fc 3, 4 and type 3 only                  |
>| status   | response status        | string        | -         | "ok"              | :heavy_check_mark:                       |
>| data     | response value         | integer array |           | [1, 0, 24, 1]     | if success                               |
>| bytes    | response byte array    | bytes array   | -         | [AB, 12, CD, ED]  | fc 3, 4 and type 2~8 only                |


#### 1.1.1 services to psmb

**bits read (FC1, FC2)**
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

**register read (FC3, FC4) - type 1, 2 (raw)**
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

**register read (FC3, FC4) - type 3 (scale)**
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

**register read (FC3, FC4) - type 4, 5 (16-bit)**
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

**register read (FC3, FC4) - type 6, 7, 8 (32-bit)**
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

#### 1.1.2 psmb to services

**bits read (FC1, FC2)**

- success:
```JavaScript
{
	"tid": 123456,
	"status": "ok",
	"data": [0,1,0,1,0,1]
}
```

- fail:
```JavaScript
{
	"tid": 123456,
	"status": "timeout"
}
```

**register read (FC3, FC4) - type 1, 2 (raw)**

- success - type 1:
```JavaScript
{
	"tid": 123456,
	"status": "ok",
	"type": 1,
	"data": [255, 1234, 789]
}
```

- success - type 2:
```JavaScript
{
	"tid": 123456,
	"status": "ok",
	"type": 2,
	"bytes": [0XFF, 0X34, 0XAB],
	"data": "112C004F12345678"
}
```

- fail:
```JavaScript
{
	"tid": 123456,
	"status": "timeout"
}
```

**register read (FC3, FC4) - type 3 (scale)**

- success - type 2:
```JavaScript
{
	"tid": 123456,
	"status": "ok",
	"bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
	"data": [22.34, 33.12, 44.56]
}
```

- fail - conversion fail:
```JavaScript
{
	"tid": 123456,
	"bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
	"status": "conversion fail"
}
```

- fail - modbus fail:
```JavaScript
{
	"tid": 123456,
	"bytes": null,
	"status": "timeout"
}
```

**register read (FC3, FC4) - type 4, 5 (16-bit)**

- success - type 4, 5:
```JavaScript
{
	"tid": 123456,
	"status": "ok",
	"bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
	"data": [255, 1234, 789]
}
```

- fail - conversion fail:
```JavaScript
{
	"tid": 123456,
	"bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
	"status": "conversion fail"
}
```

- fail - modbus fail:
```JavaScript
{
	"tid": 123456,
	"status": "timeout"
}
```

**register read (FC3, FC4) - type 6, 7, 8 (32-bit)**

- success - type 6, 7:
```JavaScript
{
	"tid": 123456,
	"status": "ok",
	"bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
	"data": [255, 1234, 789]
}
```

- success - type 8:
```JavaScript
{
	"tid": 123456,
	"status": "ok",
	"bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
	"data": [22.34, 33.12, 44.56]
}
```

- fail - conversion fail:
```JavaScript
{
	"tid": 123456,
	"bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
	"status": "conversion fail"
}
```

- fail - modbus fail:
```JavaScript
{
	"tid": 123456,
	"status": "timeout"
}
```

### 1.2 Write coil/register
Command name: **mbtcp.once.write**

>| params   | description            | type          | range          | example        | required            |
>|:---------|:-----------------------|:--------------|:---------------|:---------------|:--------------------|
>| from     | service name           | string        | -              | "web"          | optional            |
>| tid      | transaction ID         | integer       | uint64         | 12345          | :heavy_check_mark:  |
>| fc       | function code          | integer       | [1,4]          | 1              | :heavy_check_mark:  |
>| ip       | ip address             | string        | -              | 127.0.0.1      | :heavy_check_mark:  |
>| port     | port number            | string        | [1,65535]      | 502            | default: 502        |
>| slave    | slave id               | integer       | [1, 253]       | 1              | :heavy_check_mark:  |
>| addr     | register start address | integer       | -              | 23             | :heavy_check_mark:  |
>| len      | bit/register length    | integer       | -              | 20             | **FC15, 16 only**   |
>| hex      | hex/dec string flag    | bool          | [true, false]  | true           | **FC6, 16 only**    |
>| data(*)  | data to be write       | integer       | [0,1]          | 1              | **FC5 only**        |
>| data(**) | data to be write       | string        | hex/dec string | -              | **FC6, 16 only**    |
>| data(***)| data to be write       | integer array | bit array      | [1,1,0,1]      | **FC15 only**       |
>| status   | response status        | string        | -              | "ok"           | :heavy_check_mark:  |

#### 1.2.1 services to psmb

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

#### 1.2.2 psmb to services

- success:
```JavaScript
{
	"tid": 123456,
	"status": "ok"
}
```

- fail:
```JavaScript
{
	"tid": 123456,
	"status": "timeout"
}
```

### 1.3 Get TCP connection timeout
Command name: **mbtcp.timeout.read**

>| params   | description            | type          | range          | example        | required                                     |
>|:---------|:-----------------------|:--------------|:---------------|:---------------|:---------------------------------------------|
>| from     | service name           | string        | -              | "web"          | optional                                     |
>| tid      | transaction ID         | integer       | uint64         | 12345          | :heavy_check_mark:                           |
>| timeout  | timeout in usec        | integer       | [200000,~)	  | 210000         | if success                                   |
>| status   | response status        | string        | -              | "ok"           | :heavy_check_mark:                           |

#### 1.3.1 services to psmb

```JavaScript
{
    "from": "web",
    "tid": 123456
}
```

#### 1.3.2 psmb to services

- success:
```JavaScript
{
    "tid": 123456,
    "timeout": 210000,
    "status": "ok"
}
```

- fail:
```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

### 1.4 Set TCP connection timeout
Command name: **mbtcp.timeout.update**

>| params   | description            | type          | range          | example        | required                                     |
>|:---------|:-----------------------|:--------------|:---------------|:---------------|:---------------------------------------------|
>| from     | service name           | string        | -              | "web"          | optional                                     |
>| tid      | transaction ID         | integer       | uint64         | 12345          | :heavy_check_mark:                           |
>| timeout  | timeout in usec        | integer       | [200000,~)	  | 210000         | :heavy_check_mark:                           |
>| status   | response status        | string        | -              | "ok"           | :heavy_check_mark:                           |

#### 1.4.1 services to psmb

```JavaScript
{
    "from": "web",
    "tid": 123456,
    "timeout": 210000
}
```

#### 1.4.2 psmb to services

- success:
```JavaScript
{
    "tid": 123456,
    "status": "ok"
}
```

- fail:
```JavaScript
{
    "tid": 123456,
    "status": "timeout"
}
```

---

## 2. Polling requests

### 2.1 Add request
Command name: **mbtcp.poll.create**

### 2.2 Update request (interval)
Command name: **mbtcp.poll.update**

### 2.3 Read request status
Command name: **mbtcp.poll.read**

### 2.4 Delete request
Command name: **mbtcp.poll.delete**

### 2.5 Enable/Disable request
Command name: **mbtcp.poll.toggle**

### 2.6 Read all requests status
Command name: **mbtcp.polls.read**

### 2.7 Delete all requests
Command name: **mbtcp.polls.delete**

### 2.8 Enable/Disable all requests
Command name: **mbtcp.polls.toggle**

### 2.9 Import requests
Command name: **mbtcp.polls.import**

### 2.10 Export requests
Command name: **mbtcp.polls.export**

### 2.11 Read history
Command name: **mbtcp.poll.history**

---

## 3. Filter requests

### 3.1 Add filter
Command name: **mbtcp.filter.create**

### 3.2 Update filter
Command name: **mbtcp.filter.update**

### 3.3 Read filter status
Command name: **mbtcp.filter.read**

### 3.4 Delete filter
Command name: **mbtcp.filter.delete**

### 3.5 Enable/Disable filter
Command name: **mbtcp.filter.toggle**

### 3.6 Read all filters
Command name: **mbtcp.filters.read**

### 3.7 Delete all filters
Command name: **mbtcp.filters.delete**

### 3.8 Enable/Disable all filters
Command name: **mbtcp.filters.toggle**

### 3.9 Import filters
Command name: **mbtcp.filters.import**

### 3.10 Export filters
Command name: **mbtcp.filters.export**