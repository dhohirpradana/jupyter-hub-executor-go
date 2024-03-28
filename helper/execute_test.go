package helper

//
//import (
//	"fmt"
//	"github.com/gofiber/fiber/v2"
//	"gopkg.in/validator.v2"
//	"jupyter-hub-executor/entity"
//	"testing"
//)
//
//func TestExecute(t *testing.T) {
//	env, err := LoadEnv()
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	var jupyterHubExecutor *entity.JupyterHubExecutor
//
//	if err := c.BodyParser(&jupyterHubExecutor); err != nil {
//		fmt.Println(err)
//	}
//
//	if err := validator.Validate(jupyterHubExecutor); err != nil {
//		fmt.Println(err)
//	}
//
//	if c.Query("cron") == "1" {
//		jupyterHubExecutor.Cron = true
//	}
//
//	cron := &jupyterHubExecutor.Cron
//	schedulerId := &jupyterHubExecutor.SchedulerId
//	cronExpression := &jupyterHubExecutor.CronExpression
//	//user := &jupyterHubExecutor.User
//
//	if *cron {
//		if *cronExpression == "" {
//			fmt.Println(err)
//		}
//	}
//
//	token, err := GetToken()
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(token)
//
//	pbSchedulerUrl := env.PocketbaseLoginUrl
//
//	response, err := HTTPRequest(fiber.MethodGet, pbSchedulerUrl+"/"+*schedulerId, nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(response)
//}
