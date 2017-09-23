package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

var cachedClient *http.Client

func init() {
	os.Mkdir(".cache", 0755)
	// log.Println("mkdir .cache...")
	c := diskcache.New(".cache")
	t := httpcache.NewTransport(c)
	cachedClient = &http.Client{Transport: t}
}

func GetJSON(url string, v interface{}) error {
	req, err := makeRequest(url, nil, nil)
	if err != nil {
		return err
	}
	return fetch(req, v)
}

func fetch(req *http.Request, v interface{}) error {
	// log.Printf("[GET] %s\n", req.URL)
	r, err := cachedClient.Do(req)
	// log.Printf("[END] %s...\n", req.URL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return fmt.Errorf("invalid status code %d returned when fetching content from %s", r.StatusCode, req.URL)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// log.Printf("[DEBUG] %s\n", string(body))

	err = json.Unmarshal(body, v)
	return err
}

func GetJSONWithCookie(url string, cookie map[string]string, v interface{}) error {
	req, err := makeRequest(url, cookie, nil)
	if err != nil {
		return err
	}
	return fetch(req, v)
}

func makeRequest(url string, cookie map[string]string, headers map[string]string) (*http.Request, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", "curl/7.53.1")

	if cookie != nil {
		for k, v := range cookie {
			r.AddCookie(&http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	}
	if headers != nil {
		for k, v := range headers {
			r.Header.Set(k, v)
		}
	}
	return r, nil
}
