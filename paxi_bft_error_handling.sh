duration=$1
number_of_servers=$2
number_of_strugglers=$3
fast_links=$4
slow_links=$5


while true; do
    ./launch_paxibft_simulation.sh $duration $number_of_servers $number_of_strugglers $fast_links $slow_links
    wait $!
    if ! grep -q ERROR results.txt; then
        echo "No error found in results.txt"
        break
    else
        echo "Error found in results.txt"
        sleep 10
    fi
done





# Add results to the aggregated results file
echo EXPERIMENT RUN COMPLETED WITHOUT ERRORS
echo "Number of servers: $number_of_servers" >> aggregated_results.txt
echo "Number of strugglers: $number_of_strugglers" >> aggregated_results.txt
echo "Fast links: $fast_links" >> aggregated_results.txt
echo "Slow links: $slow_links" >> aggregated_results.txt
cat results.txt >> aggregated_results.txt
echo "****************************************************************" >> aggregated_results.txt
