duration=$1
number_of_servers=$2
number_of_strugglers=$3
fast_links=$4
slow_links=$5

# example usage: ./launch_paxibft_simulation.sh 60 4 2 100Kbps 10Kbps

python generate_topology.py --num_servers $number_of_servers --num_strugglers $number_of_strugglers --normal_link_capacity $fast_links --struggler_link_capacity $slow_links
cp topology.xml Kollaps/examples/paxibft/
cd PaxiBFT
# delete old docker images
docker rmi kollaps/paxibft:1.0
docker build -t kollaps/paxibft:1.0 .
cd ..

cd Kollaps/examples/
yes | ./KollapsDeploymentGenerator ./paxibft/topology.xml -s paxibft.yaml




docker stack deploy -c paxibft.yaml paxibft

echo "****************************************************************"
sleep 20
curl 127.0.0.1:8088/start
echo "Start time $(date)"
sleep $duration
curl 127.0.0.1:8088/stop
echo "End time $(date)"
echo "****************************************************************"

cd ../.. 
echo "Simulation completed"
echo "******************************** Results ********************************"
./PaxiBFT/collect_logs.sh
docker stack rm paxibft

docker stop $(docker ps -q)
docker rm $(docker ps -a -q)





