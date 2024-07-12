
image_name="kollaps/paxibft:1.0"

# The path of the log file inside the container
log_file_path_inside_container="/PaxiBFT/bin/paxi-bft.log"

# Directory on the host to save the log files
outer_log_directory="./logs"
folder_name="PaxiBFT-$(date +%Y-%m-%d_%H-%M-%S)"

host_log_directory="$outer_log_directory/$folder_name"
mkdir -p "$host_log_directory"

# Find all running containers based on the specified image name
container_ids=$(docker ps --filter "ancestor=$image_name" --format "{{.ID}}")

# Copy log files from each container to the host
for id in $container_ids; do
    container_name=$(docker inspect --format "{{.Name}}" $id | sed 's/\///')
    server_number=${container_name:0:11}
    #server_number=$(echo "$container_name" | sed -n 's/.*\(s[0-9]*\).*/\1/p')
    #echo "Copying logs from container $container_name ($id)"
    docker cp "$id:$log_file_path_inside_container" "$host_log_directory/${server_number}.log"
    docker cp "$id:/PaxiBFT/bin/config.json" "$host_log_directory/${server_number}_config.json"
done

echo "Log files have been copied to $host_log_directory"
echo "******************************** Results ********************************"
#remove results file if it exists
rm -f ./results.txt
touch ./results.txt
cat $host_log_directory/paxibft_cli.log >> ./results.txt 
rm -r $host_log_directory