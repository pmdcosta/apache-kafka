FROM anapsix/alpine-java
MAINTAINER Pedro Costa

# Kafka stable version
ENV KAFKA_VERSION 0.9.0.0
ENV SCALA_VERSION 2.11

# Download Kafka
RUN wget -q http://mirrors.fe.up.pt/pub/apache/kafka/${KAFKA_VERSION}/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz

# Install Kafka
RUN apk add --update coreutils nano
RUN tar -xzf kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz -C /opt && \
    rm kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz

# Set Kafka configuration
ENV KAFKA_HOME /opt/kafka_${SCALA_VERSION}-${KAFKA_VERSION}

# Add boot script
ADD start-kafka.sh /usr/bin/start-kafka.sh
RUN chmod +x /usr/bin/start-kafka.sh

# Add SSL configs
ADD ssl.properties /opt/kafka_${SCALA_VERSION}-${KAFKA_VERSION}/config/ssl.properties
RUN mkdir /ssl

# Set entry directory
WORKDIR /opt/kafka_${SCALA_VERSION}-${KAFKA_VERSION}

# Persistence
#VOLUME ["/kafka"]

# Entrypoint
CMD ["start-kafka.sh"]
