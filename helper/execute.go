package helper

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/validator.v2"
	"jupyter-hub-executor/entity"
	"time"
)

type JupyterHandler struct {
}

func InitJupyter() JupyterHandler {
	return JupyterHandler{}
}

func (h JupyterHandler) Execute(c *fiber.Ctx) (err error) {
	env, err := LoadEnv()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var jupyterHubExecutor *entity.JupyterHubExecutor

	if err := c.BodyParser(&jupyterHubExecutor); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := validator.Validate(jupyterHubExecutor); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if c.Query("cron") == "1" {
		jupyterHubExecutor.Cron = true
	}

	cron := &jupyterHubExecutor.Cron
	schedulerId := &jupyterHubExecutor.SchedulerId
	cronExpression := &jupyterHubExecutor.CronExpression
	//user := &jupyterHubExecutor.User

	if *cron {
		if *cronExpression == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity)
		}
	}

	token, err := GetToken()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	pbSchedulerUrl := env.PocketbaseSchedulerUrl
	jupyterUrl := env.JupyterUrl

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	url := pbSchedulerUrl + "/" + *schedulerId

	response, body, err := HTTPRequest(fiber.MethodGet, url, nil, headers)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var schedulerResponse entity.SchedulerResponse

	// Get scheduler details
	err = GetScheduler(pbSchedulerUrl, *schedulerId, &schedulerResponse)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = UnmarshalResponse(body, &schedulerResponse)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//fmt.Println("Scheduler:", schedulerResponse)

	userID := schedulerResponse.User
	pathNotebook := schedulerResponse.PathNotebook
	var user entity.User

	// Get user details
	url = env.PocketbaseUserUrl + "/" + userID

	response, body, err = HTTPRequest(fiber.MethodGet, url, nil, headers)

	err = UnmarshalResponse(body, &user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//fmt.Println("User:", user)

	port := user.JPort

	if port == 0 {
		return fiber.NewError(fiber.StatusInternalServerError, "Jupyter port is not set")
	}

	url = jupyterUrl + ":" + fmt.Sprint(port) + "/user/" + user.Username + "/api/contents"
	fmt.Println("Jupyter:", url)

	token = env.JupyterToken
	headers = map[string]string{
		"Authorization": "token " + token,
	}

	response, body, err = HTTPRequest(fiber.MethodGet, url, nil, headers)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//fmt.Println("Jupyter:", string(body))

	var directory entity.Directory

	err = UnmarshalResponse(body, &directory)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//fmt.Println("Directory Created:", directory.Created)

	// Get notebook details
	url = jupyterUrl + ":" + fmt.Sprint(port) + "/user/" + user.Username + "/api/contents/" + pathNotebook

	response, body, err = HTTPRequest(fiber.MethodGet, url, nil, headers)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//fmt.Println("Notebook:", string(body))

	var notebook entity.Notebook

	err = UnmarshalResponse(body, &notebook)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Get kernel
	now := time.Now()
	apiURL := jupyterUrl + ":" + fmt.Sprint(port) + "/user/" + user.Username + "/api"
	sessionUrl := apiURL + "/sessions?" + fmt.Sprint(now.Unix())

	kernelID, err := GetKernel(sessionUrl, pathNotebook, headers)
	fmt.Println("Kernel ID:", kernelID)

	// Execute notebook
	go func() {
		// Your notebook execution code here
		cells := notebook.Content.Cells
		jupyterWS := env.JupyterWs + ":" + fmt.Sprint(port)

		results, err := ExecuteNotebook(cells, kernelID, token, jupyterWS, apiURL)
		if err != nil {
			// Log or handle error if notebook execution fails
			fmt.Println("Notebook execution error:", err)
			return
		}

		// Process results after execution
		countOK := 0
		countError := 0
		count := len(results)

		for _, item := range results {
			if item.Status == "ok" {
				countOK++
			} else {
				countError++
			}
		}

		// Print OK, Error, Total
		fmt.Println("OK:", countOK, "Error:", countError, "Total:", count)

		// Further processing or logging can be done here

		// Example: Store results in a database, if required

	}()

	contentType := response.Header.Get("Content-Type")

	c.Set("Content-Type", contentType)

	msg := "Notebook execution initiated"

	return c.SendString(msg)
}
