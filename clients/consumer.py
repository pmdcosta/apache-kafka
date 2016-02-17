from kafka import KafkaConsumer

# To consume latest messages and auto-commit offsets
consumer = KafkaConsumer(group_id='my-group', bootstrap_servers=['172.17.0.3:9092'])
consumer.subscribe(topics=['topic'])

try:
    for message in consumer:
        print("%s:%d:%d: key=%s value=%s" % (message.topic, message.partition,
                                             message.offset, message.key,
                                             message.value))
except KeyboardInterrupt:
    consumer.close()
    print("Closed")
