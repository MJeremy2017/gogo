package scrape

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func RoundRawPrice(rp float64) float64 {
	i := float64(int64(rp))
	left := rp - i
	if left >= 0.5 {
		return i + 1
	}
	return i
}

func postAndGetJsonResponse(url string) (ticketItems, error) {
	res := ticketItems{}
	response, err := http.Post(url, "application/json", nil)
	defer response.Body.Close()
	if err != nil {
		return res, err
	}

	bytes, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
