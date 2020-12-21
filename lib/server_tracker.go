package lib

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

type PrivateServer interface {
	FetchAmount() (int, error)
	GetName() string
}

type ReturnsUserCount interface {
	UserCount() int
}

func FetchPageJsonParsed(url string, ruc ReturnsUserCount) (int, error) {
	err := FetchPageJson(url, ruc)
	if err != nil {
		return 0, err
	}
	return ruc.UserCount(), nil
}

func FetchPageJson(url string, obj interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "unable to get online users")
	}

	allData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read all data")
	}

	if err := json.Unmarshal(allData, obj); err != nil {
		return errors.Wrap(err, "unable to unmarshal content")
	}

	return nil
}

func FetchPageRegex(url string, regex *regexp.Regexp) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, errors.Wrap(err, "unable to download website")
	}

	allData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrap(err, "unable to read all data")
	}

	subMatch := regex.FindStringSubmatch(string(allData))
	if len(subMatch) == 0 {
		return 0, errors.New("unable to find submatch")
	}

	return strconv.Atoi(subMatch[1])
}
