package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"jupyter-hub-executor/entity"
)

func GetToken() (string, error) {
	env, err := LoadEnv()
	if err != nil {
		return "", err
	}

	url := env.PocketbaseLoginUrl
	authBody := entity.AuthBody{
		Identity: env.PocketbaseIdentity,
		Password: env.PocketbasePassword,
	}

	body, err := json.Marshal(authBody)
	if err != nil {
		return "", err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	_, response, err := HTTPRequest(fiber.MethodPost, url, bytes.NewReader(body), headers)
	if err != nil {
		return "", err
	}

	var adminAuth entity.AuthResponse

	err = json.Unmarshal(response, &adminAuth)
	if err != nil {
		return "", err
	}

	return adminAuth.Token, nil
}

func GetScheduler(pbSchedulerUrl string, schedulerId string, schedulerResponse *entity.SchedulerResponse) error {
	token, err := GetToken()
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	url := pbSchedulerUrl + "/" + schedulerId

	_, body, err := HTTPRequest(fiber.MethodGet, url, nil, headers)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &schedulerResponse); err != nil {
		fmt.Println("Error:", err)
	}

	err = UnmarshalResponse(body, &schedulerResponse)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}
