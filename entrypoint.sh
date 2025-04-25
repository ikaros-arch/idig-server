#!/bin/bash

# Initialize the project and add a default user if not already done
if [ ! -d "/data/K25" ]; then
  ./idig-server create K25
  ./idig-server adduser K25 admin manolis2016
fi

# Start the server
exec ./idig-server start
