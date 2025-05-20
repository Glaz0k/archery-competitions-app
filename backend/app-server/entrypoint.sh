#!/bin/sh

echo "Running migrations..."
./app-migrate
if [ $? -ne 0 ]; then
  echo "Failed to apply migrations"
  exit 1
fi

echo "Migrations applied successfully"

exec ./app-server