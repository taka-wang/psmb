#!/bin/bash
service modbusd start
/modbusd/tests/cmbserver/server &
/modbusd/tests/cmbserver/server 503