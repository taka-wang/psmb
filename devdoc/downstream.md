# Zero MQ Command Definition

# Table of contents

<!-- TOC depthFrom:2 depthTo:2 insertAnchor:false orderedList:false updateOnSave:true withLinks:true -->

- [0. Multipart message](#0-multipart-message)
- [1. Commands](#1-commands)
	- [1.1 Read request](#11-read-request)
		- [1.1.1 psmb to modbusd](#111-psmb-to-modbusd)
		- [1.1.2 modbusd to psmb](#112-modbusd-to-psmb)
	- [1.2 Write request](#12-write-request)
		- [1.2.1 psmb to modbusd](#121-psmb-to-modbusd)
		- [1.2.2 modbusd to psmb](#122-modbusd-to-psmb)
	- [1.3 Timeout request](#13-timeout-request)
		- [1.3.1 psmb to modbusd](#131-psmb-to-modbusd)
		- [1.3.2 modbusd to psmb](#132-modbusd-to-psmb)
	- [1.4 Generic fail response](#14-generic-fail-response)

<!-- /TOC -->

## 0. Multipart message

We can compose a message out of several frames, and then receiver will receive all parts of a message, or none at all.
Thanks to the all-or-nothing characteristics, we can screen what we are interested from the first frame without parsing the whole JSON payload. 
- Mode: "tcp", "rtu", others

>| Frame 1     |  Frame 2      |
>|:-----------:|:-------------:|
>| Mode        |  JSON Command |

---

## 1. Commands

>| params   | description            | type          | range     | example           |
>|:---------|:-----------------------|:--------------|:----------|:------------------|
>| tid      | transaction ID         | integer       | int64     | 12345             |
>| cmd      | command code           | integer       | -         | 1                 |
>| ip       | ip address             | string        | -         | 127.0.0.1         |
>| port     | port number            | string        | [1,65535] | 502               |
>| slave    | slave id               | integer       | [1, 253]  | 1                 |
>| addr     | register start address | integer       | -         | 23                |
>| len      | bit/register length    | integer       | -         | 20                |
>| status   | response status        | string        | -         | "ok"              |

### 1.1 Read request

#### 1.1.1 psmb to modbusd
**mbtcp read request**
```javascript
{
	"ip": "192.168.3.2",
	"port": "502",
	"slave": 22,
	"tid": 1,
	"cmd": 1,
	"addr": 250,
	"len": 10
}
```

#### 1.1.2 modbusd to psmb
**mbtcp single read response**
```javascript
{
	"tid": 1,
	"data": [1],
	"status": "ok"
}
```

**mbtcp multiple read response**
```javascript
{
	"tid": 1,
	"data": [1,2,3,4],
	"status": "ok"
}
```

### 1.2 Write request

#### 1.2.1 psmb to modbusd
**mbtcp single write request**
```javascript
{
	"ip": "192.168.3.2",
	"port": "502",
	"slave": 22,
	"tid": 1,
	"cmd": 6,
	"addr": 80,
	"data": 1234
}
```

**mbtcp multiple write request**
```javascript
{
	"ip": "192.168.3.2",
	"port": "502",
	"slave": 22,
	"tid": 3333333,
	"cmd": 16,
	"addr": 80,
	"len": 4,
	"data": [1,2,3,4]
}
```
#### 1.2.2 modbusd to psmb

**mbtcp write response**
```javascript
{
	"tid": 3333333,
	"status": "ok"
}
```

### 1.3 Timeout request

#### 1.3.1 psmb to modbusd

**mbtcp set timeout request**
```javascript
{
	"tid": 22222,
	"cmd": 50,
	"timeout": 210000
}
```

**mbtcp get timeout request**
```javascript
{
	"tid": 22222,
	"cmd": 51
}
```

#### 1.3.2 modbusd to psmb

**mbtcp set timeout response**
```javascript
{
	"tid": 22222,
	"status": "ok"
}
```

**mbtcp get timeout response**
```javascript
{
	"tid": 22222,
	"status": "ok",
	"timeout": 210000
}
```
### 1.4 Generic fail response

**mbtcp fail response**
```javascript
{
	"tid": 123456,
	"status": "fail reason"
}
```