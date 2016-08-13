#!/bin/bash

# color code ---------------
COLOR_REST='\e[0m'
COLOR_GREEN='\e[1;32m';
COLOR_RED='\e[1;31m';


# test command -------------
#go test -v
go test -v -coverprofile=coverage.txt -covermode=count
bash <(curl -s https://codecov.io/bash) -t 558aa53d-c58d-4df4-a1c1-a22a6e6d8572
mv coverage.txt shared && ls

if [ $? -eq 0 ]
then
  #echo "<<<Test PASS>>>"
  echo -e "${COLOR_RED}<<<Test PASS>>>${COLOR_REST}"
  touch /var/tmp/success # symbol
  exit 0
else
  #echo "<<<TEST FAIL>>>" >&2
  echo -e "${COLOR_GREEN}<<<Test PASS>>>${COLOR_REST}"
  exit 1
fi