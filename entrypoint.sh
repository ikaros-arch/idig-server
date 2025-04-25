#!/bin/bash

# Initialize the project and add a default user if not already done
if [ ! -d "/data/KNS-Dhaskalio" ]; then
  ./idig-server create "KNS-Dhaskalio"
  ./idig-server adduser "KNS-Dhaskalio" admin manolis2016
fi

echo "Starting the server..."

# Start the server
exec ./idig-server start
