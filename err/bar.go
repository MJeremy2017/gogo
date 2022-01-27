package err

import (
	"net/http"
	"fmt"
	"io/ioutil"
)


func DumbGetter(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("problem fetching from %s, %v", url, err)
	}

	if res.StatusCode != http.StatusOK {
        return "", BadStatusError{url, res.StatusCode}
    }

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body) // ignoring err for brevity

    return string(body), nil
}