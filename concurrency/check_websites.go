package concurrency


// a functional type, all functions with same params and return match this type
type WebsiteChecker func(string) bool

type result struct {
	url string
	res bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChan := make(chan result)

	for _, url := range urls {
		go func(u string) {
			resultChan <- result{u, wc(u)}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		r := <-resultChan
		results[r.url] = r.res
	}

	return results
}