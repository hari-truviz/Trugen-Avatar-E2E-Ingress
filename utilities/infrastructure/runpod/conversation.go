package Utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	CustomTypes "infra-ingress/customtypes"
	"io"
	"net/http"
	"os"
)

func PostJob(convReq CustomTypes.ConversationRequest) (error, CustomTypes.RequestResponse) {
	client := &http.Client{}
	convReqBytes, err := json.Marshal(convReq)
	if err != nil {
		fmt.Println("Unable to parse json")
		return err, CustomTypes.RequestResponse{}
	}
	URL := "https://api.runpod.ai/v2/" + os.Getenv("RUNPOD_ENDPOINT_ID") + "/run"
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(convReqBytes))
	if err != nil {
		fmt.Println("Error creating HTTP Request")
		return err, CustomTypes.RequestResponse{}
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("RUNPOD_API_KEY"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting response")
		return err, CustomTypes.RequestResponse{}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// Parse body to json
	response := CustomTypes.RequestResponse{}
	json.Unmarshal(body, &response)
	return nil, response
}
