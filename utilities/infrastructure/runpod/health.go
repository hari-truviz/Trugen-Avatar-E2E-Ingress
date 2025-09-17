package Utilities

import (
	"encoding/json"
	"fmt"
	CustomTypes "infra-ingress/customtypes"
	"io"
	"net/http"
	"os"
)

func IsRead() bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runpod.ai/v2/"+os.Getenv("RUNPOD_ENDPOINT_ID")+"/health", nil)
	if err != nil {
		fmt.Println("Error creating HTTP Request")
		return false
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("RUNPOD_API_KEY"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting response")
		return false
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error ready body")
		return false
	}
	var healthResponse = &CustomTypes.Health{}
	err = json.Unmarshal(body, healthResponse)
	if err != nil {
		fmt.Printf("Error parsing response: %s\n", err)
		return false
	}
	return healthResponse.Workers.Ready > 0
}
