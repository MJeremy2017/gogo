package scrape

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

func getQuantityRangeFromItem(item map[string]interface{}) string {
	q, ok := item["QuantityRange"].(string)
	if !ok {
		log.Println("failed to convert QuantityRange from item", item["QuantityRange"])
		q = ""
	}
	return q
}

func getRawPriceFromItem(item map[string]interface{}) float64 {
	p, ok := item["RawPrice"].(float64)
	if !ok {
		log.Printf("failed to convert raw price to %v float64\n", item["RawPrice"])
		p = 0.0
	}
	return p
}

func getBuyUrlFromItem(item map[string]interface{}) string {
	u, ok := item["BuyUrl"].(string)
	if !ok {
		log.Println("Failed to convert buy Url", item["BuyUrl"])
	}
	return u
}
