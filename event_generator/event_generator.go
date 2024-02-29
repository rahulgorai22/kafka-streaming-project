package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/IBM/sarama"
	"go.crwd.dev/streaming-take-home-assignment/protos"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	brokers      = flag.String("brokers", "kafka:9092", "list of brokers to connect to")
	version      = flag.String("version", sarama.DefaultVersion.String(), "Kafka cluster version")
	topic        = flag.String("topic", "cs.sensor_events", "the topic to produce messages to")
	testDataPath = flag.String("test-data-path", "/opt/test_data.json", "path to the file to load test data from")
)

func main() {
	flag.Parse()

	brokerList := strings.Split(*brokers, ",")
	if len(brokerList) == 0 {
		log.Fatal("brokers list is empty")
	}

	if len(*topic) == 0 {
		log.Fatal("no topic supplied")
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("error while creating the kafka producer: %v", err)
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error { return generateMessages(producer) })
	go func() {
		if err := group.Wait(); err != nil {
			log.Printf("received error while running: %v", err)
		}
		cancel()
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	waitForCompletion(ctx, sigterm)

	if err := producer.Close(); err != nil {
		log.Fatalf("Error closing client: %v", err)
	}
}

func waitForCompletion(ctx context.Context, sigterm chan os.Signal) {
	for {
		select {
		case <-ctx.Done():
			log.Println("context cancelled")
			return
		case <-sigterm:
			log.Println("sigterm received")
			return
		}
	}
}

func generateMessages(producer sarama.SyncProducer) error {
	testDataBytes, err := os.ReadFile(*testDataPath)
	if err != nil {
		return fmt.Errorf("failed to read test file data: %w", err)
	}

	var testData []interface{}
	err = json.Unmarshal(testDataBytes, &testData)
	if err != nil {
		return fmt.Errorf("error unmarshaling list of json data: %w", err)
	}

	for _, untypedData := range testData {
		byteContent, err := json.Marshal(untypedData)
		if err != nil {
			return fmt.Errorf("error marshaling data back into bytes: %w", err)
		}

		var sensorData protos.SensorData
		if err := protojson.Unmarshal(byteContent, &sensorData); err != nil {
			return fmt.Errorf("error marshaling bytes into proto message struct: %w", err)
		}

		sensorData.EventTimestamp = time.Now().Format(time.RFC3339)
		marshalledBytes, err := proto.Marshal(&sensorData)
		if err != nil {
			return fmt.Errorf("error marshaling sensor data protobuf to bytes: %w", err)
		}

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: *topic,
			Value: sarama.ByteEncoder(marshalledBytes),
		})

		if err != nil {
			return fmt.Errorf("error producing message to kafka topic: %w", err)
		}

		log.Printf("produced kafka message for sha256: %s, partition: %d offset: %d", sensorData.Sha256, partition, offset)

		sleepTime := rand.Intn(3000-100) + 100
		time.Sleep(time.Millisecond * time.Duration(sleepTime))
	}

	return nil
}
