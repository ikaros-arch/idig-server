#!/bin/bash

# Initialize the project and add a default user if not already done
if [ ! -d "/data/K25" ]; then
  ./idig-server create "KNS-Dhaskalio Project"
  ./idig-server adduser "KNS-Dhaskalio Project" admin manolis2016
fi

# Start the server
exec ./idig-server start
