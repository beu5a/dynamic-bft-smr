
slow_links_array=("10Kbps" "100Kbps" "200Kbps" "500Kbps" "1Mbps" "2Mbps" "3Mbps" "4Mbps" "5Mbps" "6Mbps" "7Mbps" "8Mbps" "9Mbps" "10Mbps")
fast_link="10Mbps"


for slow_link in "${slow_links_array[@]}"
do
    echo "Running benchmark with slow link: $slow_link and fast link: $fast_link | number of servers: 4 | number of strugglers: 2"
    ./paxi_bft_error_handling.sh 60 4 2 $fast_link $slow_link
    sleep 10
done