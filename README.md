# kafka-docker

Creation of custom docker images for Apache Kafka

## Contents

- Apache Zookeeper Docker Image
- Apache Kafka Docker Image
- Docker-Compose configuration
- Shell container to easily communicate with the cluster
- Python producer and consumer
- Golang producer and consumr

## Instructions

A Kafka cluster can be created using docker-compose, simply by running `docker-compose up -d`. The cluster can be scaled using `docker-compose scale kafka=5`.
If using docker-swarm, it is possible to scale the cluster across machines without changes to the images.

Alternatively, the containers can be run manually on separate machines, using the commands on the docker_example file. Simply mount the host ports using the `-p "host:container"` flag, and add the following environment variables, `KAFKA_ADVERTISED_HOST_NAME` as the hostname of the machine where the container is running, and `KAFKA_ADVERTISED_PORT` as the port mounted on the host. The env `KAFKA_ZOOKEEPER_CONNECT` should be set with the zookeeper nodes, as a comma separate list in teh format `hostname:port`. In this case, the env `HOSTNAME_COMMAND` should not be set.

The Shell container can be run using the command on the `docker_example` file. Additional control commands are also provided.

Both the python and Golang clients use the new consumer (Apache Kafka 0.9.0.0). Offsets are persisted to the cluster.
SSL authentication is not available for the python clients, but it is available on the Go clients.
