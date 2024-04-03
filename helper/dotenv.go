package helper

import (
	"github.com/joho/godotenv"
	"jupyter-hub-executor/entity"
	"os"
)

func LoadEnv() (entity.ENV, error) {
	var env entity.ENV
	err := godotenv.Load()
	if err != nil {
		return env, err
	}

	env.PocketbaseLoginUrl = os.Getenv("PB_LOGIN_URL")
	env.PocketbaseIdentity = os.Getenv("PB_MAIL")
	env.JupyterUrl = os.Getenv("JUPYTER_URL")
	env.JupyterWs = os.Getenv("JUPYTER_WS")
	env.EventUrl = os.Getenv("EVENT_URL")
	env.JupyterToken = os.Getenv("JUPYTER_TOKEN")
	env.ElasticUrl = os.Getenv("ELASTIC_URL")
	env.ElasticIndex = os.Getenv("ELASTIC_INDEX")
	env.PocketbaseUserUrl = os.Getenv("PB_USER_URL")
	env.PocketbaseSchedulerUrl = os.Getenv("PB_SCHEDULER_URL")
	env.PocketbaseNotificationUrl = os.Getenv("PB_NOTIFICATION_URL")
	env.PocketbasePassword = os.Getenv("PB_PASSWORD")

	return env, nil
}
