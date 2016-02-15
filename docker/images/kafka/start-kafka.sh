#!/bin/sh

# Optional ENV variables:
# * ADVERTISED_HOST: the external IP for the container
# * ADVERTISED_PORT: the external port for Kafka, e.g. 9092
# * BROKER_ID: the id of the broker, by default auto allocate broker ID
# * LOG_RETENTION_HOURS: the minimum age of a log file in hours to be eligible for deletion (default is 168, for 1 week)
# * LOG_RETENTION_BYTES: configure the size at which segments are pruned from the log, (default is 1073741824, for 1GB)
# * AUTO_CREATE_TOPICS: enable/disable auto creation of topics

# Set the external host and port
if [ ! -z "$HOST_NAME" ]; then
    echo "hostname: $HOST_NAME"
    sed -r -i "s/#(host.name.name)=(.*)/\1=$HOST_NAME/g" $KAFKA_HOME/config/server.properties
fi
if [ ! -z "$ADVERTISED_HOST" ]; then
    echo "advertised host: $ADVERTISED_HOST"
    sed -r -i "s/#(advertised.host.name)=(.*)/\1=$ADVERTISED_HOST/g" $KAFKA_HOME/config/server.properties
fi
if [ ! -z "$ADVERTISED_PORT" ]; then
    echo "advertised port: $ADVERTISED_PORT"
    sed -r -i "s/#(advertised.port)=(.*)/\1=$ADVERTISED_PORT/g" $KAFKA_HOME/config/server.properties
fi

# Set zookeeper host and port
if [ ! -z "$ZOOKEEPER_CONNECT" ]; then
    echo "zookeeper connect: $ZOOKEEPER_CONNECT"
    sed -r -i "s/(zookeeper.connect)=(.*)/\1=$ZOOKEEPER_CONNECT/g" $KAFKA_HOME/config/server.properties
fi

# Set Broker Id
# By default auto allocate broker ID
if [[ -z "$BROKER_ID" ]]; then
    export BROKER_ID=-1
    echo "broker id: $BROKER_ID"
    sed -r -i "s/#(zookeeper.connect)=(.*)/\1=$BROKER_ID/g" $KAFKA_HOME/config/server.properties
fi

# Allow specification of log retention policies
if [ ! -z "$LOG_RETENTION_HOURS" ]; then
    echo "log retention hours: $LOG_RETENTION_HOURS"
    sed -r -i "s/(log.retention.hours)=(.*)/\1=$LOG_RETENTION_HOURS/g" $KAFKA_HOME/config/server.properties
fi
if [ ! -z "$LOG_RETENTION_BYTES" ]; then
    echo "log retention bytes: $LOG_RETENTION_BYTES"
    sed -r -i "s/#(log.retention.bytes)=(.*)/\1=$LOG_RETENTION_BYTES/g" $KAFKA_HOME/config/server.properties
fi

# Enable/disable auto creation of topics
if [ ! -z "$AUTO_CREATE_TOPICS" ]; then
    echo "auto.create.topics.enable: $AUTO_CREATE_TOPICS"
    echo "auto.create.topics.enable=$AUTO_CREATE_TOPICS" >> $KAFKA_HOME/config/server.properties
fi

# Run Kafka
$KAFKA_HOME/bin/kafka-server-start.sh $KAFKA_HOME/config/server.properties
