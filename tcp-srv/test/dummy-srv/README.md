# Dummy upstream service

Dummy proactive service tester in golang.

## Motivation

I implement this service to test the communication among upstream services (ex. web server, websocket service and so on) and proactive service.

## Build from source code

```bash
sudo apt-get install pkg-config
curl -O https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
tar -xvf go1.6.2.linux-amd64.tar.gz
sudo mv go /usr/local
nano ~/.profile
export PATH=$PATH:/usr/local/go/bin
go get github.com/takawang/zmq3
```

## Continuous Integration

I do continuous integration with self-hosted drone.io.

---

## Test Cases

- [x] Test Timeout Operations
    - [x] `mbtcp.timeout.update` test - invalid JSON type    - (1/5)
    - [x] `mbtcp.timeout.update` test - invalid value (1)    - (2/5)
    - [x] `mbtcp.timeout.read`   test - invalid value (1)    - (3/5)
    - [x] `mbtcp.timeout.update` test - valid value (212345) - (4/5)
    - [x] `mbtcp.timeout.read`   test - valid value (212345) - (5/5)
- [x] Test One-Off Write FC5
    - [x] `mbtcp.once.write FC5` write bit test: port 502 - invalid value (2) - (1/4)
    - [x] `mbtcp.once.write FC5` write bit test: port 502 - miss from & port  - (2/4)
    - [x] `mbtcp.once.write FC5` write bit test: port 502 - valid value (0)   - (3/4)
    - [x] `mbtcp.once.write FC5` write bit test: port 502 - valid value (1)   - (4/4)
- [x] Test One-Off Write FC6
    - [x] `mbtcp.once.write FC6` write `DEC` register test: port 502 - valid value (22)         - (1/8)
    - [x] `mbtcp.once.write FC6` write `DEC` register test: port 502 - miss hex type & port     - (2/8)
    - [x] `mbtcp.once.write FC6` write `DEC` register test: port 502 - invalid value (array)    - (3/8)
    - [x] `mbtcp.once.write FC6` write `DEC` register test: port 502 - invalid hex type         - (4/8)
    - [x] `mbtcp.once.write FC6` write `HEX` register test: port 502 - valid value (ABCD)       - (5/8)
    - [x] `mbtcp.once.write FC6` write `HEX` register test: port 502 - miss port (ABCD)         - (6/8)
    - [x] `mbtcp.once.write FC6` write `HEX` register test: port 502 - invalid value (ABCD1234) - (7/8)
    - [x] `mbtcp.once.write FC6` write `HEX` register test: port 502 - invalid hex type         - (8/8)
- [x] Test One-Off Write FC15
    - [x] `mbtcp.once.write FC15` write bit test: port 502 - invalid JSON type - (1/5)
    - [x] `mbtcp.once.write FC15` write bit test: port 502 - invalid JSON type - (2/5)
    - [x] `mbtcp.once.write FC15` write bit test: port 502 - invalid value(2)  - (3/5)
    - [x] `mbtcp.once.write FC15` write bit test: port 502 - miss from & port  - (4/5)
    - [x] `mbtcp.once.write FC15` write bit test: port 502 - valid value(0)    - (5/5)
- [x] Test One-Off Write FC16
    - [x] `mbtcp.once.write FC16` write `DEC` register test: port 502 - valid value (11,22,33,44)      - (1/8)
    - [x] `mbtcp.once.write FC16` write `DEC` register test: port 502 - miss hex type & port           - (2/8)
    - [x] `mbtcp.once.write FC16` write `DEC` register test: port 502 - invalid hex type               - (3/8)
    - [x] `mbtcp.once.write FC16` write `DEC` register test: port 502 - invalid length                 - (4/8)
    - [x] `mbtcp.once.write FC16` write `HEX` register test: port 502 - valid value (ABCD1234)         - (5/8)
    - [x] `mbtcp.once.write FC16` write `HEX` register test: port 502 - miss port (ABCD)               - (6/8)
    - [x] `mbtcp.once.write FC16` write `HEX` register test: port 502 - invalid hex type (11,22,33,44) - (7/8)
    - [x] `mbtcp.once.write FC16` write `HEX` register test: port 502 - invalid length                 - (8/8)
- [x] Test One-Off Read FC1
    - [x] `FC1` read bits test: port 502 - length 1
    - [x] `FC1` read bits test: port 502 - length 7
    - [x] `FC1` read bits test: port 502 - Illegal data address
    - [x] `FC1` read bits test: port 503 - length 7
- [x] Test One-Off Read FC2
    - [x] `FC2` read bits test: port 502 - length 1
    - [x] `FC2` read bits test: port 502 - length 7
    - [x] `FC2` read bits test: port 502 - Illegal data address
    - [x] `FC2` read bits test: port 503 - length 7
- [x] Test One-Off Read FC3
    - [x] `FC3` read bytes Type 1 test: port 502
    - [x] `FC3` read bytes Type 2 test: port 502
    - [x] `FC3` read bytes Type 3 length 4 test: port 502
    - [x] `FC3` read bytes Type 3 length 7 test: port 502 - invalid length
    - [x] `FC3` read bytes Type 4 length 4 test: port 502 - Order: AB
    - [x] `FC3` read bytes Type 4 length 4 test: port 502 - Order: BA
    - [x] `FC3` read bytes Type 4 length 4 test: port 502 - miss order
    - [x] `FC3` read bytes Type 5 length 4 test: port 502 - Order: AB
    - [x] `FC3` read bytes Type 5 length 4 test: port 502 - Order: BA
    - [x] `FC3` read bytes Type 5 length 4 test: port 502 - miss order
    - [x] `FC3` read bytes Type 6 length 8 test: port 502 - Order: AB
    - [x] `FC3` read bytes Type 6 length 8 test: port 502 - Order: BA
    - [x] `FC3` read bytes Type 6 length 8 test: port 502 - miss order
    - [x] `FC3` read bytes Type 6 length 7 test: port 502 - invalid length
    - [x] `FC3` read bytes Type 7 length 8 test: port 502 - Order: AB
    - [x] `FC3` read bytes Type 7 length 8 test: port 502 - Order: BA
    - [x] `FC3` read bytes Type 7 length 8 test: port 502 - miss order
    - [x] `FC3` read bytes Type 7 length 7 test: port 502 - invalid length
    - [x] `FC3` read bytes Type 8 length 8 test: port 502 - order: ABCD
    - [x] `FC3` read bytes Type 8 length 8 test: port 502 - order: DCBA
    - [x] `FC3` read bytes Type 8 length 8 test: port 502 - order: BADC
    - [x] `FC3` read bytes Type 8 length 8 test: port 502 - order: CDAB
    - [x] `FC3` read bytes Type 8 length 7 test: port 502 - invalid length
    - [x] `FC3` read bytes: port 502 - invalid type
- [x] Test One-Off Read FC4
    - [x] `FC4` read bytes Type 1 test: port 502
    - [x] `FC4` read bytes Type 2 test: port 502
    - [x] `FC4` read bytes Type 3 length 4 test: port 502
    - [x] `FC4` read bytes Type 3 length 7 test: port 502 - invalid length
    - [x] `FC4` read bytes Type 4 length 4 test: port 502 - Order: AB
    - [x] `FC4` read bytes Type 4 length 4 test: port 502 - Order: BA
    - [x] `FC4` read bytes Type 4 length 4 test: port 502 - miss order
    - [x] `FC4` read bytes Type 5 length 4 test: port 502 - Order: AB
    - [x] `FC4` read bytes Type 5 length 4 test: port 502 - Order: BA
    - [x] `FC4` read bytes Type 5 length 4 test: port 502 - miss order
    - [x] `FC4` read bytes Type 6 length 8 test: port 502 - Order: AB
    - [x] `FC4` read bytes Type 6 length 8 test: port 502 - Order: BA
    - [x] `FC4` read bytes Type 6 length 8 test: port 502 - miss order
    - [x] `FC4` read bytes Type 6 length 7 test: port 502 - invalid length
    - [x] `FC4` read bytes Type 7 length 8 test: port 502 - Order: AB
    - [x] `FC4` read bytes Type 7 length 8 test: port 502 - Order: BA
    - [x] `FC4` read bytes Type 7 length 8 test: port 502 - miss order
    - [x] `FC4` read bytes Type 7 length 7 test: port 502 - invalid length
    - [x] `FC4` read bytes Type 8 length 8 test: port 502 - order: ABCD
    - [x] `FC4` read bytes Type 8 length 8 test: port 502 - order: DCBA
    - [x] `FC4` read bytes Type 8 length 8 test: port 502 - order: BADC
    - [x] `FC4` read bytes Type 8 length 8 test: port 502 - order: CDAB
    - [x] `FC4` read bytes Type 8 length 7 test: port 502 - invalid length
    - [x] `FC4` read bytes: port 502 - invalid type
- [x] Test Poll Single Requests
    - [x] `mbtcp.poll.create FC1` read bits test: port 503 - miss name
    - [x] `mbtcp.poll.create/mbtcp.poll.delete FC1` read bits test: port 503 - interval 1
    - [x] `mbtcp.poll.update/mbtcp.poll.delete FC1` read bits test: port 503 - miss name
    - [x] `mbtcp.poll.update/mbtcp.poll.delete FC1` read bits test: port 503 - interval 2
    - [x] `mbtcp.poll.read/mbtcp.poll.delete FC1` read bits test: port 503 - miss name
    - [x] `mbtcp.poll.read/mbtcp.poll.delete FC1` read bits test: port 503
    - [x] `mbtcp.poll.read/mbtcp.poll.toggle/mbtcp.poll.delete FC1` read bits test: port 503
    - [x] `mbtcp.poll.toggle/mbtcp.poll.read/mbtcp.poll.delete FC1` read bits test: port 503 - enable
- [x] Test Polls Requests
    - [x] `mbtcp.polls.read FC1` read 2 poll reqeusts
    - [x] `mbtcp.polls.read/mbtcp.polls.read/mbtcp.polls.delete FC1` read 50 poll requests
    - [x] `mbtcp.polls.read/mbtcp.polls.read/mbtcp.polls.toggle/mbtcp.polls.delete FC1` read 20 poll requests


---

## License

MIT