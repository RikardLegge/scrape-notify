package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadToken(t *testing.T) {
	assert(t, os.Unsetenv("PUSHBULLET_TOKEN"))

	_, err := loadToken()
	assertFalse(t, err, "Should not be able to load token when non is set")
	assert(t, os.Setenv("PUSHBULLET_TOKEN", "value"))

	token, err := loadToken()
	assert(t, err, "Should be able to load token when set")
	assert(t, token == "value")
}

func TestSearchPageCsb(t *testing.T) {
	csbData, err := ioutil.ReadFile("./test_data/csb.html")
	entries, err := searchPage(csb, string(csbData))
	assert(t, err)
	assert(t, len(entries) == 1)
}

func TestSearchPage(t *testing.T) {
	page := "a: 1, b: 2; a: 3, b: 4"
	scraper := Scraper{
		keys:  []string{"a", "b"},
		value: "[^\\w]*(\\d)[^\\w]*",
	}
	entries, err := searchPage(scraper, page)
	assert(t, err)
	assert(t, len(entries) == 2)
}
