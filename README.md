# psmb

[![GoDoc](https://godoc.org/github.com/taka-wang/psmb?status.svg)](http://godoc.org/github.com/taka-wang/psmb)
[![Go Report Card](https://goreportcard.com/badge/github.com/taka-wang/psmb)](https://goreportcard.com/report/github.com/taka-wang/psmb)
[![Build Status](http://drone.cmwang.net/api/badges/taka-wang/psmb/status.svg)](http://drone.cmwang.net/taka-wang/psmb)
[![CircleCI](https://circleci.com/gh/taka-wang/psmb.svg?style=shield)](https://circleci.com/gh/taka-wang/psmb)
[![codecov](https://codecov.io/gh/taka-wang/psmb/branch/master/graph/badge.svg)](https://codecov.io/gh/taka-wang/psmb)

Proactive service library for [modbusd](https://github.com/taka-wang/modbusd)

---

## Environment variables

> Why environment variable? Refer to the [12 factors](http://12factor.net/)

- CONF_PSMBTCP: config file path
- EP_BACKEND: endpoint of remote service discovery server (optional)

## Contracts (Interfaces)

- IProactiveService: proactive service
- IReaderTaskDataStore:  read/poll task data store
- IWriterTaskDataStore: write task data store
- IHistoryDataStore: history data store
- IFilterDataStore: filter data store
- IConfig: config management

## Golang package management

- I adopted [glide](https://glide.sh/) as package management system for this repository.

## Worker Pool Model

### Request from upstream

![uml](http://uml.cmwang.net:8000/plantuml/svg/5Sh13O0W3030LNG0QSBJZxDKQ908XPGsnEtLzzsQEHIBP5AMIxMF7K1mkfJrijC6IMYinEf2gw1uupQH4tIh1IeE9O58lRIdVWdCH_VJuLy0)

### Response from downstream

![uml](http://uml.cmwang.net:8000/plantuml/svg/5Sh13O0W3030LNG0QVJfnragD4AaG4eRulRg-svEnMSBP9AdgDhw3Y0ut9KqsccTnDUYLDJvog1uupAmND2CCp1s9O50BTU7lmHXC_VJiRu0)

---

## Up and Running

### Docker Compose

```bash
docker-compose build --pull
docker-compose up
```

### Manually

#### Install libzmq (3.2.5)

```bash
wget https://github.com/zeromq/zeromq3-x/releases/download/v3.2.5/zeromq-3.2.5.tar.gz
tar xvzf zeromq-3.2.5.tar.gz
cd zeromq-3.2.5
./configure
make
sudo make install
sudo ldconfig
```

#### Build

```bash
# install glide
curl https://glide.sh/get | sh

# install dependencies
glide install

# build
cd tcp-srv
go build -o psmbtcp-srv
```

---

## Continuous Integration

I do continuous integration and deploy docker images after git push by self-hosted drone.io server, [circleci](https://circleci.com/) service, [codecov](https://codecov.io/) service and [dockerhub]((https://hub.docker.com/r/edgepro/c-modbus-slave/)) service.

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

![uml](http://uml.cmwang.net:8000/plantuml/svg/5SZ13O0W3030LNG0QVpwSPPI2H1R8BGDwnllUNjjnFuadxmLiw4NmCGLShNYqJLDwirIiq1TmF35Os7BC5mO1DNI169KXQ4ImzytdXy0)