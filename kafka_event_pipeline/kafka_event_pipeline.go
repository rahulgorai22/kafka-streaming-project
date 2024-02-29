package main

import (
	"context"
	"flag"
	"github.com/IBM/sarama"
	"github.com/golang/protobuf/proto"
	"go.crwd.dev/streaming-take-home-assignment/kafka_event_pipeline/utils"
	"go.crwd.dev/streaming-take-home-assignment/protos" // Adjust this import path to where your generated protobuf Go package is.
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	kafkaBrokers = flag.String("kafka-brokers", "kafka:9092", "list of kafka brokers to connect to")
)

var logger = utils.GetLogger()

func main() {
	flag.Parse()
	// ***************************
	//Set the kafka configurations
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_0_0_0
	// ***************************

	// Create a new consumer group
	group, err := sarama.NewConsumerGroup(strings.Split(*kafkaBrokers, ","), "sensor_processor_group", config)
	if err != nil {
		logger.Sugar().Fatalf("Error in creating consumer group due to error: %v", err)
	}
	consumer := Consumer{
		ready: make(chan bool),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			if err := group.Consume(ctx, []string{utils.KAFKA_TOPIC}, &consumer); err != nil {
				logger.Sugar().Fatalf("Received an error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	logger.Sugar().Infof("Sarama consumer up and running !! %v", consumer)

	// Wait for termination signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	logger.Sugar().Infof("Termination signal received. Shutting down !! %v", consumer)
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	cu := utils.CassandraUtilities{}
	dbSession, err := cu.CreateDatabaseSession()
	if err != nil {
		logger.Sugar().Errorf("Failed to create database session due to error %v", err)
	}
	defer cu.CloseSession(dbSession)

	for message := range claim.Messages() {
		var data protos.SensorData
		if err := proto.Unmarshal(message.Value, &data); err != nil {
			logger.Sugar().Errorf("Failed to unmarshal message: %v", err)
		}
		// Log the message timestamp and topic
		logger.Sugar().Debugf("Message claimed: timestamp = %v, topic = %s\n", message.Timestamp, message.Topic)

		sensorData := sensorData{
			Platform: data.Platform,
			SHA256:   data.Sha256,
		}

		APIClient := utils.APIClient{}
		classification, score, err := APIClient.Execute(sensorData.SHA256, sensorData.Platform.String(), data.CustomerId)
		if err != nil {
			logger.Sugar().Errorf("Customer [%s]: Error in making api call", data.CustomerId)
		}

		err = cu.InsertData(dbSession, data.Sha256, score, classification, data.EventTimestamp)
		if err != nil {
			logger.Sugar().Errorf("Customer [%s]: Errror in db call due to %v\n", data.CustomerId, err)
		} else {
			logger.Sugar().Infof("Customer [%s]: Successfully updated database", data.CustomerId)
		}
		session.MarkMessage(message, "")
	}

	return nil
}
