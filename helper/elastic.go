package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"io"
	"jupyter-hub-executor/entity"
	"time"
)

type ESScheduler struct {
	Path        string                `json:"path"`
	UserId      string                `json:"uid"`
	SchedulerId string                `json:"scheduler_id"`
	CellResults []entity.ESCellResult `json:"cell_results"`
	Date        string                `json:"date"`
	DateFinish  string                `json:"date_finish"`
	Ok          int                   `json:"success"`
	Error       int                   `json:"error"`
	Executed    int                   `json:"executed"`
	Total       int                   `json:"total"`
	ElapsedTime time.Duration         `json:"elapsed_time"`
}

func getESClient() (*elasticsearch.Client, error) {
	env, err := LoadEnv()
	if err != nil {
		fmt.Println(err.Error())
	}

	//fmt.Println(env.ElasticUrl)

	cfg := elasticsearch.Config{
		Addresses: []string{env.ElasticUrl},
	}
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		return nil, err
	}

	return es, nil
}

func (a ESScheduler) StoreToES() {
	//fmt.Println("A:", a)

	env, err := LoadEnv()
	if err != nil {
		fmt.Println(err.Error())
	}

	es, err := getESClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	start := time.Now()
	finishTime := start.UTC().Format("2006-01-02T15:04:05.999999Z")

	a.DateFinish = finishTime

	body, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err.Error())
	}

	//fmt.Println(env.ElasticIndex)

	res, err := es.Index(
		env.ElasticIndex,
		bytes.NewReader(body),
		es.Index.WithDocumentID(fmt.Sprintf("%d", time.Now().UnixNano())),
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
