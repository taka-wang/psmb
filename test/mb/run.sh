#!/bin/bash
service modbusd start
/modbusd/tests/cmbserver/server "127.0.0.1" 502 &
/modbusd/tests/cmbserver/server "127.0.0.1" 503