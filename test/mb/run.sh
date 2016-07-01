#!/bin/bash
service modbusd start
/modbusd/tests/cmbserver/server &
/modbusd/tests/cmbserver/server 127.0.0.1 503