# Run Zookeeper
docker run -d --name zookeeper --hostname zookeeper pmdcosta/zookeeper

# Run kafka
docker run -it --rm --name kafka_1 --hostname kafka_1 --link zookeeper -e HOST_NAME=kafka_1 -e ADVERTISED_PORT=9092 -e ZOOKEEPER_CONNECT=zookeeper:2181 pmdcosta/kafka bash

# Possible Java networking error
docker run -it --rm --link zookeeper -e ZOOKEEPER_CONNECT=zookeeper:2181 pmdcosta/kafka bash
