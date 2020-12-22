package lib

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

func FetchPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get online users")
	}
	defer resp.Body.Close()

	allData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read all data")
	}

	return allData, nil
}

func FetchPageJson(url string, obj interface{}) error {
	allData, err := FetchPage(url)
	if err != nil {
		return errors.Wrap(err, "unable to fetch page")
	}

	if err := json.Unmarshal(allData, obj); err != nil {
		return errors.Wrap(err, "unable to unmarshal content")
	}

	return nil
}

func FetchPageRegex(url string, regex *regexp.Regexp) (int, error) {
	allData, err := FetchPage(url)
	if err != nil {
		return 0, errors.Wrap(err, "unable to fetch page")
	}

	subMatch := regex.FindStringSubmatch(string(allData))
	if len(subMatch) == 0 {
		return 0, errors.New("unable to find submatch")
	}

	return strconv.Atoi(subMatch[1])
}
