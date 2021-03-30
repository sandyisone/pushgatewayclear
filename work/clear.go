package work

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func Work(baseUrl string, ttl time.Duration) {

	baseUrl = strings.TrimSuffix(baseUrl, "/")

	for {
		log.Println("begin clear work...")
		clear(baseUrl, ttl);
		log.Println("end clear work")
		time.Sleep(time.Minute * 1)

	}
}


func clear(baseUrl string, ttl time.Duration) error {

	metricUrl := fmt.Sprintf("%s/api/v1/metrics", baseUrl)
	response, err := http.Get(metricUrl)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	r := new(Response)
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}

	//过期截止时间
	expiration := time.Now().Add(0 - ttl)
	//fmt.Println(expiration.Format("2006-01-02 15:04:05"))

	for _, m := range r.Data {

		//fmt.Println(m.PushTimeSeconds.Timestamp.Format("2006-01-02 15:04:05"))
		//过期
		expired := m.PushTimeSeconds.Timestamp.Before(expiration)

		if !expired {
			continue
		}

		labels := m.PushTimeSeconds.Metrics[0].Labels
		job, ok := labels["job"]
		if !ok {
			continue
		}

		//http://xxx/metrics/job@base64/5LqR5Y2X5a6d5Y2O5bCP5a2m
		deleteUrl := fmt.Sprintf("%s/metrics/job@base64/%s", baseUrl,
			base64.RawURLEncoding.EncodeToString([]byte(job)))

		for k, v := range m.PushTimeSeconds.Metrics[0].Labels {
			if k != "" && v != "" && k != "job" {
				deleteUrl = deleteUrl + "/" + k + "/" + v
			}
		}

		err = delete(deleteUrl);
		if err != nil {
			log.Printf("delete job: %s (%s) error: %s ", job, deleteUrl, err.Error())
		}else{
			log.Printf("delete job: %s (%s) ", job, deleteUrl)
		}

	}
	return nil
}

func delete(deleteUrl string) error {

	//deleteUrl = url.QueryEscape(deleteUrl)
	//fmt.Println(deleteUrl)

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", deleteUrl, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("http status: %v", resp.Status)
	}
	return nil
}


type Response struct {
	Status    string `json:"status"`
	Data      []Data `json:"data,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
	Error     string `json:"error,omitempty"`
}

type Data struct {
	PushTimeSeconds TimestampGauge `json:"push_time_seconds"`
}

type TimestampGauge struct {
	Timestamp time.Time `json:"time_stamp"`
	Metrics   []Metric  `json:"metrics"`
}

type Metric struct {
	Labels map[string]string `json:"labels"`
}
