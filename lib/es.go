package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type ElasticSearch struct {
	URL   string `yaml:"url"`
	Index string `yaml:"index"`
}

type ElasticSearchDocument struct {
	Timestamp                time.Time `json:"@timestamp"`
	ServerName               string    `json:"server_name"`
	OnlineCount              int       `json:"online_count"`
	ResponseTimeMilliseconds int       `json:"response_time_ms"`
}

func (es *ElasticSearch) Write(doc *ElasticSearchDocument) error {
	b, err := json.Marshal(doc)
	if err != nil {
		return errors.Wrap(err, "unable to marshal document")
	}

	resp, err := http.Post(
		fmt.Sprintf(
			"%s/%s-%d.%d.%d/_doc/",
			es.URL,
			es.Index,
			doc.Timestamp.Year(),
			doc.Timestamp.Month(),
			doc.Timestamp.Day(),
		),
		"application/json",
		bytes.NewReader(b),
	)

	if err != nil {
		return errors.Wrap(err, "unable to send data")
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return errors.Errorf("got error code %d", resp.StatusCode)
	}

	return nil
}
