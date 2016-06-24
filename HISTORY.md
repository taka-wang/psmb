# Done
- [x] evaluate scheduler
- [x] evaluate scheduler with zmq
- [x] hacking gocron package
- [x] setup southbound docker env for testing
- [x] integrate with travis ci

# TODO
- [ ] define request type

# Backlog:
- if scheduler is stop, trigger requests directly.
- handle default port if not set
- handle data length
    - read: default = 1
    - write: check length with data, or set it automatically
- if the length of the response data equal to 1, should we put it data 'array'
- check timeout interval range
- check polling interval range

