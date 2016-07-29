# Done

- [x] evaluate scheduler
- [x] evaluate scheduler with zmq
- [x] hacking gocron package
- [x] setup southbound docker env for testing
- [x] integrate with travis ci
- [x] implement binary module unit test
- [x] implement byte array to uint16 array conversion
- [x] implement byte array to int16 array conversion
- [x] implement byte array to uint32 array conversion
- [x] implement byte array to int32 array conversion
- [x] implement byte array to float32 array conversion
- [x] implement register array to byte array conversion
- [x] implement byte array to hex string conversion
- [x] implement heximal string to register array conversion
- [x] implement decimal string to register array converson
- [x] support x86_64, arm docker CI
- [x] implement linear scaling conversion
- [x] bit string to uint16 array conversion
- [x] finish downstream devdoc
- [x] implement upstream struct test cases
- [x] implement downstream struct test cases
- [x] unify types definitions
- [x] support docker compose
- [x] refactor main flow
- [x] implement uint test for mbtcp.once.write
- [x] implement integration test for mbtcp.once.read
- [x] implement integration test for mbtcp.once.write
- [x] handle default port and fc15/fc16 length
- [x] define polling commands
- [x] define MbtcpCmdType for modbusd
- [x] implement mutex lock for write map
- [x] implement integration test for mbtcp.timeout.read
- [x] implement integration test for mbtcp.timeout.update
- [x] implement NaiveResponser
- [x] Refactor read task map mechanism
- [x] Refactor write task to OO
- [x] implement factory pattern
- [x] Generalize proactive service implementation
- [x] Re-org test cases and integrate with drone.io
- [x] Handle default port if not set
- [x] implement GetByTID and GetByName for ReadTaskType
- [x] implement all poll request handlers (except mbtcp.poll.history)
- [x] implement all polls request handlers

# Backlog:

- if scheduler is stop, trigger requests directly.
- sqlite3 for history
- task test cases
- table size limit
- fix handleResponse last return issue
- add to log and filter

# Module

- post processing (event mapper)
    - filter
    - on changed
- database?
