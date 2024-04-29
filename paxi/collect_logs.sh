
image_name="kollaps/paxi:1.0"

# The path of the log file inside the container
log_file_path_inside_container="/paxi/bin/paxi.log"

# Directory on the host to save the log files
outer_log_directory="./logs"
folder_name="Folder-$(date +%Y-%m-%d_%H-%M-%S)"

host_log_directory="$outer_log_directory/$folder_name"
mkdir -p "$host_log_directory"

# Find all running containers based on the specified image name
container_ids=$(docker ps --filter "ancestor=$image_name" --format "{{.ID}}")

# Copy log files from each container to the host
for id in $container_ids; do
    container_name=$(docker inspect --format "{{.Name}}" $id | sed 's/\///')
    server_number=${container_name:0:11}
    #server_number=$(echo "$container_name" | sed -n 's/.*\(s[0-9]*\).*/\1/p')
    echo "Copying logs from container $container_name ($id)"
    docker cp "$id:$log_file_path_inside_container" "$host_log_directory/${server_number}.log"
done

echo "Log files have been copied to $host_log_directory"