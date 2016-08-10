# psmb

[![Build Status](http://drone.cmwang.net/api/badges/taka-wang/psmb/status.svg)](http://drone.cmwang.net/taka-wang/psmb)
[![GoDoc](https://godoc.org/github.com/taka-wang/psmb?status.svg)](http://godoc.org/github.com/taka-wang/psmb)

Proactive service library for [modbusd](https://github.com/taka-wang/modbusd)

---

## Environment variables

- CONF_PSMBTCP: config file location
- EP_BACKEND: remote service discovery endpoint (optional)

## Contracts (Interfaces)

- IProactiveService: proactive service
- IReaderTaskDataStore:  read/poll task map
- IWriterTaskDataStore: write task map
- IHistoryDataStore: history map
- IFilterDataStore: filter map
- IConfig: config

## Worker Pool

### Request

![uml](http://uml.cmwang.net:8000/plantuml/svg/5Sh13O0W3030LNG0QUBJRIeqOAI5b3R4xTNttNf9h9x8jIx5h8y3G766V5dnqmTfR68a5c9ZCBnncoWPkdC5nc6aaAZNzl2NmFSKVde1)


### Response

![uml](http://uml.cmwang.net:8000/plantuml/svg/5Sh13O0W3030LNG0QUBJRIeqOAI5b3R4xTNttNf9h9x8jIx5h8y3G766V5dnqmTfR68a5c9ZCBnncoWPkdC5nc6aaAZNzl2N8EqUVde1)


## Docker Compose

```bash
docker-compose build --pull
docker-compose up
```

## Continuous Integration

I do continuous integration and build docker images after git push by self-hosted drone.io server and [dockerhub]((https://hub.docker.com/r/takawang/c-modbus-slave/)) service.

## Deployment Diagram

![uml](http://uml.cmwang.net:8000/plantuml/svg/5Sh13O0W3030LNG0QVJfnraAD42aGA0DSNlrVRUcuh9wqfwNADB62T1ncf0agjL1tTKYLCIuoY1uupQn16ZA6HY7K0TFBTU7lmHji3M_NVln0W00)

---

## Unit tests

- binary
- types

## Test cases

### Binary

- [x] Bytes to 16-bit integer array tests
    - [x] `BytesToUInt16s` in big endian order - (1/4)
    - [x] `BytesToUInt16s` in little endian order - (2/4)
    - [x] `BytesToInt16s` in big endian order - (3/4)
    - [x] `BytesToInt16s` in little endian order - (4/4)
- [x] Bytes to 32-bit integer array tests
    - [x] `BytesToUInt32s` in (ABCD) Big Endian order - (1/4)
    - [x] `BytesToUInt32s` in (DCBA) Little Endian order - (2/4)
    - [x] `BytesToUInt32s` in (BADC) Mid-Big Endian order - (3/4)
    - [x] `BytesToUInt32s` in (CDAB) Mid-Little Endian order - (4/4)
    - [x] `BytesToInt32s` in (ABCD) Big Endian order - (1/4)
    - [x] `BytesToInt32s` in (DCBA) Little Endian order - (2/4)
    - [x] `BytesToInt32s` in (BADC) Mid-Big Endian order - (3/4)
    - [x] `BytesToInt32s` in (CDAB) Mid-Little Endian order - (4/4)
- [x] Bytes to 32-bit float array tests
    - [x] `BytesToFloat32s` in (ABCD) Big Endian order - (1/4)
    - [x] `BytesToFloat32s` in (DCBA) Little Endian order - (2/4)
    - [x] `BytesToFloat32s` in (BADC) Mid-Big Endian order - (3/4)
    - [x] `BytesToFloat32s` in (CDAB) Mid-Little Endian order - (4/4)
- [x] Bytes/registers utility tests
    - [x] `BitStringToUInt8s` test
    - [x] `BitStringToUInt8s` test - left comma
    - [x] `BitStringToUInt8s` test - right comma
    - [x] `BitStringToUInt8s` test - left, right comma
    - [x] `RegistersToBytes` test
    - [x] `BytesToHexString` test
    - [x] `DecimalStringToRegisters` test
    - [x] `DecimalStringToRegisters` test - left comma
    - [x] `DecimalStringToRegisters` test - right comma
    - [x] `DecimalStringToRegisters` test - left, right comma
    - [x] `HexStringToRegisters` test
    - [x] `HexStringToRegisters` test - wrong length
    - [x] `LinearScalingRegisters` test
    - [x] `LinearScalingRegisters` test - (0,0,0,0)
    - [x] `LinearScalingRegisters` test - reverse

### Types

#### Upstream structure test

- [x] One-off modbus tcp struct tests
    - [x] `mbtcp.once.read` request test
    - [x] `mbtcp.once.read` response test
- [x] get/set modbus tcp timeout struct tests
    - [x] `mbtcp.timeout.read` request test
    - [x] `mbtcp.timeout.read` response test
    - [x] `mbtcp.timeout.update` request test
    - [x] `mbtcp.timeout.update` response test
    - [x] `mbtcp.once.write` request test
    - [x] `mbtcp.once.write` response test

#### Downstream structure test

- [x] modbus tcp downstreamstruct tests
    - [x] `read` request test
    - [x] `single read` response test
    - [x] `multiple read` response test
    - [x] `single write` request test
    - [x] `multiple write` request test
    - [x] `set timeout` request test
    - [x] `get timeout` response test

---

## UML

![uml](http://uml.cmwang.net:8000/plantuml/svg/5SZB3O0W303GLNG0gMSlpb8gGI8KqazONtt7jnQcwbTogSjjDlG049mX5xizkYQXpfRO0lK6XWzk4pd3y5QXeLeIe8ggCBJ5yFUvVru0)
