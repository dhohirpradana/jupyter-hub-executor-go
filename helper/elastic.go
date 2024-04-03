package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
	"io"
	"jupyter-hub-executor/entity"
	"os"
	"time"
)

type ESScheduler struct {
	Path        string `json:"path"`
	UserId      string `json:"uid"`
	SchedulerId string `json:"scheduler_id"`
	CellResults []entity.ESCellResult
	Date        string `json:"date"`
}

func getESClient() (*elasticsearch.Client, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	elasticUrl := os.Getenv("ELASTIC_URL")

	cfg := elasticsearch.Config{
		Addresses: []string{elasticUrl},
	}
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		return nil, err
	}

	return es, nil
}

func (a *ESScheduler) StoreToES() {
	fmt.Println("A:", a)

	env, err := LoadEnv()
	if err != nil {
		fmt.Println(err.Error())
	}

	es, err := getESClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	t := time.Now()
	formattedTime := t.UTC().Format("2006-01-02T15:04:05.999999Z")
	a.Date = formattedTime

	body, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err.Error())
	}

	res, err := es.Index(
		env.ElasticIndex,
		bytes.NewReader(body),
		es.Index.WithDocumentID(fmt.Sprintf("%d", t.UnixNano())),
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.IsError() {
		fmt.Println(err.Error())
	}

	fmt.Println("Success store data to ES")
}
