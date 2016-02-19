package main

import (
	"github.com/Shopify/sarama"
	"gopkg.in/bsm/sarama-cluster.v2"

	"flag"
	"strings"
	"os"
	"os/signal"
	"syscall"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
	"log"
)

var (
	brokerList 	= flag.String("brokerList", "172.17.0.3:9092", "The Kafka brokers to connect to, as a comma separated list")
	topicList	= flag.String("topicList", "topic", "The topics to publish the message to, as a comma separated list")
	groupId 	= flag.String("group", "go", "The shared consumer group name")
	offset 		= flag.String("offset", "newest", "The offset to start with. Can be 'oldest' or 'newest'")

	certFile 	= flag.String("certificate", "", "The optional certificate file for client authentication")
	keyFile 	= flag.String("key", "", "The optional key file for client authentication")
	caFile 		= flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	verifySsl 	= flag.Bool("verify", false, "Optional verify ssl certificates chain")

	logger = log.New(os.Stderr, "", log.LstdFlags)
)


func main() {
	// Parse command line flags
	flag.Parse()

	// Ensure flags were supplied
	if *brokerList == "" || *topicList == "" || *groupId == "" || *offset == ""{
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Separate lists
	brokers := strings.Split(*brokerList, ",")
	topics := strings.Split(*topicList, ",")
	logger.Printf("Kafka Brokers:\t %v\n", brokers)
	logger.Printf("Kafka Topics:\t %v\n", topics)
	logger.Printf("Kafka Group:\t %v\n", *groupId)

	// Create new consumer
	consumer := newKafkaConsumer(brokers, topics, *groupId, *offset, *certFile, *keyFile, *caFile, *verifySsl)

	// Clean shutdown
	defer closeConsumer(*consumer)

	// Receive Messages
	go receiveMessage(*consumer)

	// Wait for signal
	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-wait
}

// Create Kafka Consumer
func newKafkaConsumer(brokers []string, topics []string, group string, offset string, caPath string, certPath string, keyPath string, verify bool) *cluster.Consumer {
	// Create configurations
	config := cluster.NewConfig()

	// Select starting offset
	switch offset {
	case "oldest":
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "newest":
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		logger.Printf("Wrong offset flag, should be 'oldest' or 'newest'\n")
		os.Exit(64)
	}

	// Retention of commit logs
	config.Consumer.Offsets.CommitInterval = 1

	// Setup SSL Configuration
	sslConfig := sslConfig(caPath, certPath, keyPath, verify)
	if sslConfig != nil {
		config.Net.TLS.Config = sslConfig
		config.Net.TLS.Enable = true
	}

	// Create consumer
	consumer, err := cluster.NewConsumer(brokers, group, topics, config)
	if err != nil {
		logger.Printf("There was an error building the Producer: %s\n", err)
	}

	return consumer
}

// Create SSL configuration
func sslConfig(caPath string, certPath string, keyPath string, verify bool) *tls.Config {
	// Check if credentials were provided
	if caPath == "" && certPath == "" && keyPath == "" {
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

// Receive Messages
func receiveMessage(consumer cluster.Consumer) {
	logger.Printf("Receiving Messages:\n")
	for msg := range consumer.Messages() {
		logger.Printf("[%s/%d/%d]\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		consumer.MarkOffset(msg, "")
	}
}

// Close Consumer
func closeConsumer(consumer cluster.Consumer)  {
	logger.Printf("Shuting down...\n")
	if err := consumer.Close(); err != nil {
		logger.Printf("Error closing Consumer: %s\n", err)
	}
}































