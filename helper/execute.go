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

	var jupyterHubExecutor entity.JupyterHubExecutor

	if err := c.BodyParser(&jupyterHubExecutor); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := validator.Validate(jupyterHubExecutor); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if c.Query("cron") == "1" {
		jupyterHubExecutor.Cron = true
	}

	cron := jupyterHubExecutor.Cron
	schedulerId := jupyterHubExecutor.SchedulerId
	cronExpression := jupyterHubExecutor.CronExpression
	//user := &jupyterHubExecutor.User

	if cron {
		if cronExpression == "" {
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

	url := pbSchedulerUrl + "/" + schedulerId

	response, body, err := HTTPRequest(fiber.MethodGet, url, nil, headers)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var schedulerResponse entity.SchedulerResponse

	// Get scheduler details
	err = GetScheduler(pbSchedulerUrl, schedulerId, &schedulerResponse)
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

	url = jupyterUrl + ":" + fmt.Sprint(port) + "/user/jupyter/api/contents"
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
	url = jupyterUrl + ":" + fmt.Sprint(port) + "/user/jupyter/api/contents/" + pathNotebook

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
	apiURL := jupyterUrl + ":" + fmt.Sprint(port) + "/user/jupyter/api"
	sessionUrl := apiURL + "/sessions?" + fmt.Sprint(now.Unix())

	kernelID, err := GetKernel(sessionUrl, pathNotebook, headers)
	fmt.Println("Kernel ID:", kernelID)

	// Execute notebook
	cells := notebook.Content.Cells
	jupyterWS := env.JupyterWs + ":" + fmt.Sprint(port)

	results := &[]entity.CellResult{}

	UpdateSchedulerStatus(pbSchedulerUrl, schedulerId, "running")

	contentType := response.Header.Get("Content-Type")

	c.Set("Content-Type", contentType)

	// Process results after execution
	countOK := 0
	countError := 0
	count := len(cells)

	var esCellResults []entity.ESCellResult

	go func() {
		start := time.Now()
		err = ExecuteNotebook(cells, kernelID, token, jupyterWS, apiURL, pbSchedulerUrl, schedulerId, results)
		if err != nil {
			fmt.Println("Notebook execution error:", err)
			return
		}

		for _, result := range *results {
			if result.Status == "ok" {
				countOK++
			} else {
				countError++
			}

			esCellResults = append(esCellResults, entity.ESCellResult{
				Cell:     result.Cell,
				CellType: result.CellType,
				Status:   result.Status,
				Message:  result.Message,
			})
		}

		if countError == 0 {
			UpdateSchedulerStatus(pbSchedulerUrl, schedulerId, "success")
		} else {
			UpdateSchedulerStatus(pbSchedulerUrl, schedulerId, "failed")
		}

		elapsed := time.Since(start)
		fmt.Println("OK:", countOK, "Error:", countError, "Executed:", countOK+countError, "Total:", count, "Execution time: %s\n", elapsed)

		var esScheduler ESScheduler

		esScheduler.SchedulerId = schedulerId
		esScheduler.Path = pathNotebook
		esScheduler.UserId = userID
		esScheduler.CellResults = esCellResults
		esScheduler.Ok = countOK
		esScheduler.Error = countError
		esScheduler.Executed = countOK + countError
		esScheduler.Total = count
		esScheduler.ElapsedTime = elapsed

		//fmt.Println(esScheduler)

		esScheduler.StoreToES()
	}()

	return c.SendString("Notebook execution initiated.")
}
