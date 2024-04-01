package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"jupyter-hub-executor/entity"
	"log"
	"net/http"
	"sync"
	"time"
)

func ExecuteWS(cellSource, kernel, token, jupyterWS, apiURL string) (map[string]interface{}, error) {
	uuid4 := uuid.New()
	msgID := uuid.New()
	now := time.Now()
	formattedDate := now.Format("2006-01-02T15:04:05.999Z")

	uri := fmt.Sprintf("%s/user/jupyter/api/kernels/%s/channels?session_id=%s&token=%s", jupyterWS, kernel, uuid4, token)
	//fmt.Println("uri", uri)

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

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	// Creating websocket connection
	ws, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)

	// Sending message
	err = ws.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		return nil, err
	}

	for {
		// Always connect to the websocket
		if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
			return nil, err
		}

		// Receiving response
		_, response, err := ws.ReadMessage()
		if err != nil {
			return nil, err
		}

		var responseJSON map[string]interface{}
		if err := json.Unmarshal(response, &responseJSON); err != nil {
			return nil, err
		}

		content := responseJSON["content"].(map[string]interface{})
		msgState := responseJSON["header"].(map[string]interface{})["msg_type"].(string)

		if msgState == "input_request" {
			restartKernel(kernel, apiURL, token)
			return map[string]interface{}{"status": "error", "msg": "input prompt"}, nil
		}

		if msgState == "error" {
			errMsg := content["traceback"].(string)
			return map[string]interface{}{"status": "error", "msg": errMsg}, nil
		}

		if status, ok := content["status"].(string); ok {
			if status == "error" {
				errMsg := content["traceback"].(string)
				return map[string]interface{}{"status": "error", "msg": errMsg}, nil
			}
			return map[string]interface{}{"status": status, "msg": "Success"}, nil
		}
	}
}

func ExecuteNotebook(cells []entity.CodeCell, kernelID, token, jupyterWS, apiURL string, results *[]entity.CellResult) error {
	var wg sync.WaitGroup
	//var results []entity.CellResult

	for index, cell := range cells {
		wg.Add(1)
		cellSource := cell.Source
		cellType := cell.CellType

		go func(i int, cell entity.CodeCell) {
			defer wg.Done()
			if cellType == "code" && cellSource != "" {
				res, err := ExecuteWS(cellSource, kernelID, token, jupyterWS, apiURL)
				if err != nil {
					*results = append(*results, entity.CellResult{
						Cell:      i + 1,
						CellType:  cellType,
						CellValue: cellSource,
						Status:    "error",
						Message:   err.Error(),
					})
				} else {
					*results = append(*results, entity.CellResult{
						Cell:       i + 1,
						CellType:   cellType,
						CellValue:  cellSource,
						Status:     res["status"].(string),
						Message:    res["msg"].(string),
						Additional: res,
					})
				}
			} else {
				*results = append(*results, entity.CellResult{
					Cell:      i + 1,
					CellType:  cellType,
					CellValue: cellSource,
					Status:    "ok",
					Message:   "Success",
				})
			}
		}(index, cell)
	}

	wg.Wait()

	return nil
}

func GetKernel(apiURL string, pathNotebook string, headers map[string]string) (string, error) {
	now := time.Now()
	sessionUrl := apiURL + "/sessions?" + fmt.Sprint(now.Unix())

	_, body, err := HTTPRequest(fiber.MethodGet, sessionUrl, nil, headers)
	if err != nil {
		return "", err
	}

	var sessions []entity.SessionResponse

	err = UnmarshalResponse(body, &sessions)
	if err != nil {
		return "", err
	}

	if len(sessions) == 0 {
		return "", errors.New("no kernels found")
	}

	var kernelID string
	for _, item := range sessions {
		if item.Path == pathNotebook {
			kernelID = item.Kernel.ID
			break
		}
	}

	return kernelID, nil
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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

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
