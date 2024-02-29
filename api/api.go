package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
)

var platformTypes = map[string]bool{
	"PLATFORM_WINDOWS": true,
	"PLATFORM_OSX":     true,
	"PLATFORM_LINUX":   true,
}

type RequestBody struct {
	Data         string `json:"data"`
	PlatformType string `json:"platform_type"`
}

func main() {
	http.HandleFunc("/classify", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Error reading request body: %v", err)
			return
		}

		errMsg := make(map[string]string)
		var requestBody RequestBody
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errMsg["error"] = "unable to parse parameters!"

			jsonResp, err := json.Marshal(errMsg)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal: %v", err)
			}
			w.Write(jsonResp)
			return
		}

		if len(requestBody.Data) == 0 || len(requestBody.PlatformType) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			if len(requestBody.Data) == 0 {
				errMsg["error"] = "data not populated!"
			} else {
				errMsg["error"] = "platform_type field not populated!"
			}

			jsonResp, err := json.Marshal(errMsg)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal: %v", err)
			}
			w.Write(jsonResp)
			return
		}

		if ok := platformTypes[requestBody.PlatformType]; !ok {
			w.WriteHeader(http.StatusInternalServerError)
			errMsg["error"] = "invalid platform type!"

			jsonResp, err := json.Marshal(errMsg)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal: %v", err)
			}
			w.Write(jsonResp)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(requestBody.Data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errMsg["error"] = "data parameter not a valid base64 string!"

			jsonResp, err := json.Marshal(errMsg)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal: %v", err)
			}
			w.Write(jsonResp)
			return
		}

		bits := binary.LittleEndian.Uint64(decoded)
		float := math.Float64frombits(bits)
		for float > 1 {
			float = float / 10
		}

		classification := ""
		if float > 0.5 {
			classification = "malicious"
		} else {
			classification = "benign"
		}

		respMsg := map[string]interface{}{
			"classification": classification,
			"score":          float,
		}

		w.WriteHeader(http.StatusOK)
		jsonResp, err := json.Marshal(respMsg)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal: %v", err)
		}
		w.Write(jsonResp)
		return

	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
