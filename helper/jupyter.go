package helper

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func executeWS(index int, cellSource, kernel, token, jupyterWS, apiURL, user string) (map[string]interface{}, error) {
	// Generate UUIDs
	uuid4 := uuid.New()
	msgID := uuid.New()

	// Get current formatted date
	formattedDate := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")

	// WebSocket URI
	uri := fmt.Sprintf("%s/user/jupyter/api/kernels/%s/channels?session_id=%s&token=%s", jupyterWS, kernel, uuid4, token)
	log.Println("uri:", uri)

	// Prepare message
	message := map[string]interface{}{
		"header": map[string]interface{}{
			"date":     formattedDate,
			"msg_id":   msgID,
			"msg_type": "execute_request",
			"session":  uuid4,
			"username": "",
			"version":  "5.2",
		},
		"parent_header": map[string]interface{}{},
		"metadata": map[string]interface{}{
			"editable":     true,
			"slideshow":    map[string]interface{}{"slide_type": ""},
			"tags":         []interface{}{},
			"trusted":      true,
			"deletedCells": []interface{}{},
			"recordTiming": false,
		},
		"content": map[string]interface{}{
			"code":             cellSource,
			"silent":           false,
			"store_history":    true,
			"user_expressions": map[string]interface{}{},
			"allow_stdin":      true,
			"stop_on_error":    true,
		},
		"buffers": []interface{}{},
	}

	// Convert message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	log.Println("message:", string(messageJSON))

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Send message
	err = conn.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		return nil, err
	}

	// Wait for response
	_, response, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// Process response
	var responseJSON map[string]interface{}
	err = json.Unmarshal(response, &responseJSON)
	if err != nil {
		return nil, err
	}

	// Handle response based on message type
	msgState := responseJSON["header"].(map[string]interface{})["msg_type"].(string)
	content := responseJSON["content"].(map[string]interface{})
	switch msgState {
	case "input_request":
		// Restart kernel and return error message
		restartKernel(kernel, apiURL, token)
		return map[string]interface{}{"status": "error", "msg": "input prompt"}, nil
	case "error":
		// Return error message
		errMsg := content["traceback"].(string)
		return map[string]interface{}{"status": "error", "msg": errMsg}, nil
	default:
		// Return status and message
		status := content["status"].(string)
		return map[string]interface{}{"status": status, "msg": content["traceback"].(string)}, nil
	}
}

func restartKernel(kernel, apiURL, token string) {
	// Get current time
	now := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")

	// Construct URL
	url := fmt.Sprintf("%s/kernels/%s/restart?%s", apiURL, kernel, now)

	// Create HTTP client
	client := &http.Client{}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal("Error creating HTTP request:", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending HTTP request:", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: %s", resp.Status)
	}

	// Read response body
	// Note: If response body contains useful information, you can handle it here
	// Otherwise, you may omit this part
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	//     log.Fatal("Error reading response body:", err)
	// }

	// Print response
	// fmt.Println("Restart kernel:", string(body))
}
