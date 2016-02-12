FROM anapsix/alpine-java
MAINTAINER Pedro Costa

# Zookeeper Stable version
ENV ZOOKEEPER_VERSION 3.4.6

# Download Zookeeper
RUN wget -q http://mirror.vorboss.net/apache/zookeeper/zookeeper-${ZOOKEEPER_VERSION}/zookeeper-${ZOOKEEPER_VERSION}.tar.gz && \
    wget -q https://www.apache.org/dist/zookeeper/KEYS && \
    wget -q https://www.apache.org/dist/zookeeper/zookeeper-${ZOOKEEPER_VERSION}/zookeeper-${ZOOKEEPER_VERSION}.tar.gz.asc && \
    wget -q https://www.apache.org/dist/zookeeper/zookeeper-${ZOOKEEPER_VERSION}/zookeeper-${ZOOKEEPER_VERSION}.tar.gz.md5

# Install Zookeeper
RUN tar -xzf zookeeper-${ZOOKEEPER_VERSION}.tar.gz -C /opt

# Set Zookeeper Configuration
RUN mv /opt/zookeeper-${ZOOKEEPER_VERSION}/conf/zoo_sample.cfg /opt/zookeeper-${ZOOKEEPER_VERSION}/conf/zoo.cfg
ENV ZK_HOME /opt/zookeeper-${ZOOKEEPER_VERSION}
RUN sed  -i "s|/tmp/zookeeper|$ZK_HOME/data|g" $ZK_HOME/conf/zoo.cfg; mkdir $ZK_HOME/data

# Add boot script
ADD start-zk.sh /usr/bin/start-zk.sh
RUN chmod +x /usr/bin/start-zk.sh

# Set entry directory
WORKDIR /opt/zookeeper-${ZOOKEEPER_VERSION}

# Persistence
#VOLUME ["/opt/zookeeper-${ZOOKEEPER_VERSION}/conf", "/opt/zookeeper-${ZOOKEEPER_VERSION}/data"]

# Entrypoint
EXPOSE  2181 2888 3888
#CMD /usr/sbin/sshd && bash /usr/bin/start-zk.sh
CMD bash /usr/bin/start-zk.sh