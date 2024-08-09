package main

import (
	"crypto/tls"
	fmevent "hola2-mv-consumer/fm-event"
	mvinterface "hola2-mv-consumer/mv-interface"
	mvnode "hola2-mv-consumer/mv-node"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	kafkaBrokers := os.Getenv("KAFKA_SERVER")
	esURL := os.Getenv("ES_SERVER")
	esUser := os.Getenv("ES_USER")
	esPassword := os.Getenv("ES_PASSWORD")

	consumer, err := sarama.NewConsumer([]string{kafkaBrokers}, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %s", err)
	}
	defer consumer.Close()

	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
		Username:  esUser,
		Password:  esPassword,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	var wg sync.WaitGroup

	topics := []string{"MV-NODE-5MIN-BB", "MV-NODE-1MIN-BB", "MV-NODE-5MIN-BH", "MV-NODE-1MIN-BH", "MV-INTERFACE-5MIN-BB", "MV-INTERFACE-1MIN-BB", "MV-INTERFACE-5MIN-BH", "MV-INTERFACE-1MIN-BH", "FM-EVENT-BB", "FM-EVENT-BH"}

	for _, topic := range topics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			switch {
			case strings.HasPrefix(topic, "MV-NODE"):
				mvnode.ProcessTopicMvNode(consumer, es, topic)
			case strings.HasPrefix(topic, "MV-INTERFACE"):
				mvinterface.ProcessTopicMvInterface(consumer, es, topic)
			case strings.HasPrefix(topic, "FM-EVENT"):
				fmevent.ProcessTopicFmEvent(consumer, es, topic)
			default:
				log.Printf("No handler for topic: %s", topic)
			}
		}(topic)
	}
	wg.Wait()
}
