while true; do
    if  grep -q ERROR results.txt; then
        echo "Error found in results.txt"
        sleep 10
    else
        echo "No error found in results.txt"
        break
    fi
done