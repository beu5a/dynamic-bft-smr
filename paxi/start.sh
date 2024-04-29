#!/bin/bash

role=$1
id=$2

# Before starting, delete all previous logs
if [[ $id == "1.1" ]]; then
    wait 5
elif [[ $id == "1.2" ]]; then
    wait 10
elif [[ $id == "1.3" ]]; then
    wait 15
elif [[ $id == "1.4" ]]; then
    wait 20
fi

if [[ $role == "server" ]]; then
    echo "Server role"
    echo "ID: $id"
    ./server -log_dir=. -log_level=info -id $id
elif [[ $role == "client" ]]; then
    echo "Client role"
    wait 30
    ./client -id $id -config config.json
else
    echo "Invalid role"
fi
