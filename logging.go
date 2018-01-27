package main

import (
    "fmt"
    "log"
    "os"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/Shopify/sarama"
)

// write to local or remote kafka instance
func send_kafka(s []byte) bool {
        config := sarama.NewConfig()
        config.Producer.RequiredAcks = sarama.WaitForAll
        config.Producer.Retry.Max = 5

        brokers := []string{"localhost:9092"}
        config.Producer.Return.Successes = true
        producer, err := sarama.NewSyncProducer(brokers, config)
        if err != nil {
                panic(err)
        }

        defer func() {
                if err := producer.Close(); err != nil {
                        panic(err)
                }
        }()

        topic := "s3_alerts"
        msg := &sarama.ProducerMessage{
                Topic: topic,
                Value: sarama.StringEncoder(s),
        }

        partition, offset, err := producer.SendMessage(msg)
        if err != nil {
                panic(err)
                return false
        }

        fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
        return true
}

// write to local log file
func write_logfile(l *s3.ListObjectsOutput) {
    f, err := os.OpenFile("s3_results.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("error opening file: %v", err)
    }
    defer f.Close()

    log.SetOutput(f)
    log.Println(l)
}
