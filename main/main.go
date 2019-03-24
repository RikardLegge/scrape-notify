package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Push struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
}

type Scraper struct {
	title string
	url   string
	value string
	keys  []string
}

var csb = Scraper{
	title: "Nytt lägenhetsförråd tillgängligt",
	url:   "https://www.chalmersstudentbostader.se/bo-hos-oss/lagenhetsforrad",
	value: ".*?>([^<]*).*\n?.*",
	keys:  []string{"Förrådsnr", "Adress", "Storlek", "Hyra", "Typ", "Ledigt from"},
}

func main() {
	token, err := loadToken()
	if err != nil {
		panic(err)
	}

	for _, s := range []Scraper{csb} {
		entries, err := scrape(s)
		if err != nil {
			panic(err)
		}
		for _, e := range entries {
			if hasSent(e) {
				continue
			}
			err = cacheSend(e)
			if err != nil {
				panic(err)
			}
			err = push(Push{
				Type:  "link",
				Title: s.title,
				Body:  e,
				URL:   s.url,
			}, token)
			if err != nil {
				panic(err)
			}
		}
	}
}

func cacheSend(body string) error {
	f, err := os.OpenFile(".cache", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(body + "\n"))
	if err != nil {
		return err
	}
	return f.Close()
}

func hasSent(body string) bool {
	cacheData, err := ioutil.ReadFile(".cache")
	if err != nil {
		return false
	}
	cache := string(cacheData)
	return strings.Contains(body, cache)
}

func scrape(ex Scraper) ([]string, error) {
	html, err := fetchWebsite(ex.url)
	if err != nil {
		return nil, err
	}
	pairList, err := searchPage(ex, html)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, pair := range pairList {
		var messages []string
		for key, val := range pair {
			messages = append(messages, key+": "+val)
		}
		results = append(results, strings.Join(messages, "\n"))
	}
	return results, nil
}

func loadToken() (string, error) {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == "PUSHBULLET_TOKEN" {
			return pair[1], nil
		}
	}
	return "", fmt.Errorf("PUSHBULLET_TOKEN not found in the environment")
}

func searchPage(ex Scraper, html string) ([]map[string]string, error) {
	value := ex.value
	keys := ex.keys
	regex := ""
	for _, key := range keys {
		regex += key + value
	}
	extractRE, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	trimRE := regexp.MustCompile("\\s+")

	groups := extractRE.FindAllStringSubmatch(html, -1)
	hits := make([]map[string]string, 0)
	for _, values := range groups {
		pairs := make(map[string]string, 0)
		for i, key := range keys {
			value := string(trimRE.ReplaceAll([]byte(values[i+1]), []byte(" ")))
			value = strings.TrimSpace(value)
			pairs[key] = value
		}
		hits = append(hits, pairs)
	}
	return hits, nil
}

func fetchWebsite(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func push(toPush Push, token string) error {
	data, err := json.Marshal(toPush)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", "https://api.pushbullet.com/v2/pushes", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Access-Token", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(body))
	}
	return nil
}
