#!/bin/bash

# Load .env variables and export them
if [ -f /tmp/.env ]; then
    echo "Loading environment variables from .env"
    export $(grep -v '^#' /tmp/.env | xargs)
else
    echo "Error: .env file not found!"
    exit 1
fi
