package utils

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type APIClient struct {
	// Extend the API Client struct to embed more client operations.
}

var (
	apiHost = flag.String("api-host", "http://api:8080", "api host to connect to")
)

func (ac *APIClient) Execute(platformData, platformType, customerId string) (string, float64, error) {
	flag.Parse()
	// Construct the URL for the API call
	url := fmt.Sprintf("%s/classify", *apiHost)
	// Construct the payload
	payload := map[string]interface{}{
		"data":          platformData,
		"platform_type": platformType,
	}
	// Marshal the sensor data into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", 0, fmt.Errorf("customer [%s]: Failed to marshal sensor data: %v", customerId, err)
	}

	// Make the HTTP POST request to the API
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return "", 0, fmt.Errorf("customer [%s]: Failed to call API: %v", customerId, err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("customer [%s]: Failed to read response body: %v", customerId, err)
	}

	// Parse the response
	var response struct {
		Classification string  `json:"classification"`
		Score          float64 `json:"score"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", 0, fmt.Errorf("customer [%s]: failed to parse API response: %v", customerId, err)
	}
	return response.Classification, response.Score, nil
}
