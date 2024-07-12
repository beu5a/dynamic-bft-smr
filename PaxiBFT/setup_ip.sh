#!/bin/bash

# Filename for temporary storing IP addresses and ports
temp_file="hosts.temp"

# Clear the file if it exists or create a new one
echo "" > $temp_file

# Fetch IP addresses for each service and write to the temp file
for i in $(seq 1 4); do
    # Assuming hostname format s$i-$KOLLAPS_UUID, e.g., s1-uuid123
    j=$(($i - 1))
    ip=$(host paxibft_s$j-$KOLLAPS_UUID | grep -E -o "([0-9]{1,3}[\.]){3}[0-9]{1,3}")
    
    # Write each line in the format needed for the JSON keys and values
    echo "\"1.$i\": \"tcp://$ip:1735\"" >> $temp_file
    echo "\"1.$i\": \"http://$ip:9080\"" >> $temp_file
done

echo "IP addresses fetched and written to $temp_file"
