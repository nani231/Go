package kafka

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

const (
	kafkaConn     = "localhost:9092"
	producerTopic = "output-data"
	filepath      = "E://temp//SAP.txt"
)

func main() {
	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	file, err := os.Open(filepath)
	msg, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(err)
	}

	publish(string(msg), producer)

	producer.Close() // To close broker
}

func initProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

func publish(message string, producer sarama.SyncProducer) {
	msg := &sarama.ProducerMessage{
		Topic: producerTopic,
		Value: sarama.StringEncoder(message),
	}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println(p, o)
		fmt.Println("Error publish: ", err.Error())
	}

}
