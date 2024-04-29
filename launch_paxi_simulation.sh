duration=$1

cd paxi
docker build -t kollaps/paxi:1.0 .
cd ..

cd Kollaps/examples/
./KollapsDeploymentGenerator ./paxi/topology.xml -s paxi.yaml




docker stack deploy -c paxi.yaml paxi


sleep 20
curl 127.0.0.1:8088/start
echo "Start time $(date)"
sleep $duration
curl 127.0.0.1:8088/stop
echo "End time $(date)"


cd ../.. 
./paxi/collect_logs.sh
docker stack rm paxi



