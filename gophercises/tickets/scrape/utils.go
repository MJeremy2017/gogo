package scrape

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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

func getStringFromMap(item map[string]interface{}, key string) string {
	q, ok := item[key].(string)
	if !ok {
		log.Println("failed to convert string from item", item[key])
		q = ""
	}
	return q
}

func getFloatFromMap(item map[string]interface{}) float64 {
	p, ok := item["RawPrice"].(float64)
	if !ok {
		log.Printf("failed to convert raw price to %v float64\n", item["RawPrice"])
		p = 0.0
	}
	return p
}

func SortEventTicketsByPrice(events []Event) []Event {
	var res []Event
	for _, e := range events {
		sort.Slice(e.Tickets, func(i, j int) bool {
			return e.Tickets[i].Price < e.Tickets[j].Price
		})
		res = append(res, e)
	}
	return res
}

func SaveEventsToJson(events []Event) {
	f, err := os.Create("scrape/event.json")
	defer f.Close()
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(events)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
}

func LoadJsonToEvents(fp string) ([]Event, error) {
	var events []Event
	f, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// returns format in 2023-02-19T15:00:00
func formatDateHour(dt, hour string) string {
	ms := map[string]string{
		"Jan": "01",
		"Feb": "02",
		"Mar": "03",
		"Apr": "04",
		"May": "05",
		"Jun": "06",
		"Jul": "07",
		"Aug": "08",
		"Sep": "09",
		"Oct": "10",
		"Nov": "11",
		"Dec": "12",
	}
	s := strings.Split(dt, " ")
	day, mon := s[0], ms[s[1]]
	year, _, _ := time.Now().Date()
	return strconv.Itoa(year) + "-" + mon + "-" + day + "T" + hour + ":00"
}

func formatDollarSignPrice(price string) float64 {
	price = strings.TrimSpace(price)
	var pStr string
	f := false
	for _, ch := range price {
		if f {
			pStr += string(ch)
		}
		if string(ch) == "$" {
			f = true
		}
	}

	p, err := strconv.Atoi(pStr)
	if err != nil {
		log.Println("failed to convert price", pStr)
		return 0.0
	}
	return float64(p)
}
