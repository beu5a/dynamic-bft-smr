#!/bin/bash

# Path to the .env file on the host
ENV_FILE_PATH="./.env"

# Check if .env file exists
if [ ! -f "$ENV_FILE_PATH" ]; then
    echo "Error: .env file does not exist at $ENV_FILE_PATH"
    exit 1
fi

# Loop through all running Docker containers
docker ps --format '{{.Names}}' | grep -E 's[0-9]+|client' | while read container_name; do
    echo "Processing $container_name"

    # Copy .env file to the container
    docker cp "$ENV_FILE_PATH" "$container_name:/tmp/.env"
    echo "$container_name updated with .env variables."
done

echo "All matching containers have been processed."
