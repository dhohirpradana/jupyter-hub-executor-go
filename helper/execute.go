package helper

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/validator.v2"
	"jupyter-hub-executor/entity"
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

	token, err := TokenGet()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	pbSchedulerUrl := env.PocketbaseSchedulerUrl

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
	err = SchedulerGet(pbSchedulerUrl, *schedulerId, &schedulerResponse)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = UnmarshalResponse(body, &schedulerResponse)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fmt.Println("Scheduler:", schedulerResponse)

	contentType := response.Header.Get("Content-Type")

	c.Set("Content-Type", contentType)

	return c.Send(body)
}
