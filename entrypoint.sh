#!/bin/bash

# Initialize the project and add a default user if not already done
if [ ! -d "/data/$PROJECT_NAME" ]; then
  ./idig-server create "$PROJECT_NAME"
  ./idig-server adduser "$PROJECT_NAME" "$ADMIN_USER" "$ADMIN_PASSWORD"
fi

echo "Starting the server..."

# Start the server
exec ./idig-server start
