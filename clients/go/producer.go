package main

import (
	"github.com/Shopify/sarama"

	"flag"
	"os"
	"strings"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
	"log"
)

var (
	brokerList 	= flag.String("brokerList", "172.17.0.3:9093", "The Kafka brokers to connect to, as a comma separated list")
	topic		= flag.String("topic", "topic", "The topic to publish the message to")

	certFile 	= flag.String("certificate", "/home/pmdcosta/DockerProjects/kafka-docker/ssl/signed-req.pem", "The optional certificate file for client authentication")
	keyFile 	= flag.String("key", "/home/pmdcosta/DockerProjects/kafka-docker/ssl/key.pem", "The optional key file for client authentication")
	caFile 		= flag.String("ca", "/home/pmdcosta/DockerProjects/kafka-docker/ssl/ca-cert", "The optional certificate authority file for TLS client authentication")
	verifySsl 	= flag.Bool("verify", true, "Optional verify ssl certificates chain")

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {
	// Parse command line flags
	flag.Parse()

	// Ensure flags were supplied
	if *brokerList == "" || *topic == ""{
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Separate the supplied brokers
	brokers := strings.Split(*brokerList, ",")
	logger.Printf("Kafka Brokers:\t %v\n", brokers)

	// Create new producer
	producer := newKafkaProducer(brokers, *caFile, *certFile, *keyFile, *verifySsl)

	// Close Producer
	defer closeProducer(producer)

	message := "Ohayo Sekai"

	// Send Message
	sendMessage(producer, *topic, message)
}

// Create Kafka Producer
func newKafkaProducer(brokers []string, caPath string, certPath string, keyPath string, verify bool) sarama.SyncProducer {
	// Create configurations
	config := sarama.NewConfig()

	// Wait for all in-sync replicas to ack the message
	config.Producer.RequiredAcks = sarama.WaitForAll

	// Setup SSL Configuration
	sslConfig := sslConfig(caPath, certPath, keyPath, verify)
	if sslConfig != nil {
		config.Net.TLS.Config = sslConfig
		config.Net.TLS.Enable = true
	}

	// Create producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Printf("There was an error building the Producer: %s\n", err)
	}

	return producer
}

// Create SSL configuration
func sslConfig(caPath string, certPath string, keyPath string, verify bool) *tls.Config {
	// Check if credentials were provided
	if caPath == "" || certPath == "" || keyPath == "" {
		logger.Printf("No SSL config provided\n")
		return nil
	}

	// Load CA certificate
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		logger.Printf("Error loading ca: %s\n", err)
	}

	// Load certificate and key
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		logger.Printf("Error loading key and cert: %s\n", err)
	}

	// Create CA certificate pool
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(ca)

	// Load SSL to config
	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caPool,
		InsecureSkipVerify: verify,
	}

	return config
}

// Publish message
func sendMessage(producer sarama.SyncProducer, topic string, payload string) {
	message := sarama.ProducerMessage{
		// The Kafka topic for this message
		Topic: topic,
		// The actual message to store in Kafka
		Value: sarama.StringEncoder(payload),
		// No message key, so messages will be distributed randomly over partitions
	}

	// Send Message
	partition, offset, err := producer.SendMessage(&message)
	if err != nil {
		logger.Printf("Error sending data: %s\n", err)
	} else {
		logger.Printf("[%s/%d/%d] Message successfully published\n", topic, partition, offset)
	}

}

// Close Producer
func closeProducer(producer sarama.SyncProducer)  {
	err := producer.Close()
	if err != nil {
		logger.Printf("Error closing Producer: %s\n", err)
	}
}