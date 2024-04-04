package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func HTTPRequest(method string, url string, body io.Reader, headers map[string]string) (*http.Response, []byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(resp.Body)

	if err != nil {
		return nil, nil, err
	}

	if resp == nil {
		return nil, nil, errors.New("response is nil")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, responseBody, nil
}

func UnmarshalResponse(body []byte, v interface{}) error {
	err := json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

//func GetNextRun(cronExpr string, tTime time.Time) (time.Time, error) {
//	specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
//	schedule, err := specParser.Parse(cronExpr)
//	if err != nil {
//		return time.Now(), err
//	}
//
//	nextRunTime := schedule.Next(tTime)
//
//	return nextRunTime, nil
//}
