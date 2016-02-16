from kafka import KafkaProducer
from kafka.common import KafkaError

producer = KafkaProducer(bootstrap_servers=['172.17.0.3:9092'])

# Asynchronous by default
future = producer.send('topic', b'Hidden Message')

# Block for 'synchronous' sends
try:
    record_metadata = future.get(timeout=10)
except KafkaError:
    print("Error")
    pass

# Successful result returns assigned partition and offset
print(record_metadata.topic)
print(record_metadata.partition)
print(record_metadata.offset)
