package main

import "go.crwd.dev/streaming-take-home-assignment/protos"

type Consumer struct {
	ready chan bool
}

type sensorData struct {
	Platform protos.Platform `json:"platform_type"`
	SHA256   string          `json:"data"`
	// Add other fields as needed
}
