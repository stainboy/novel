package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type toc struct {
	Chapterlist struct {
		Chapters []struct {
			ChapterId   int    `json:"chapterId"`
			RrderNum    int    `json:"orderNum"`
			ChapterName string `json:"chapterName"`
		} `json:"chapters"`
		PageSize int `json:"pageSize"`
		PageNum  int `json:"pageNum"`
	} `json:"chapterlist"`
}

func Test_fetch(t *testing.T) {
	var toc toc
	req, err := makeRequest("http://m.zongheng.com/h5/ajax/chapter/list?h5=1&bookId=342974&pageNum=1&pageSize=2500&chapterId=0&asc=0", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fetch2(req, &toc)

	log.Print("ok")
	log.Print(len(toc.Chapterlist.Chapters))
	log.Print("o2k")
}

func fetch2(req *http.Request, v interface{}) error {
	log.Printf("[GET] %s\n", req.URL)
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

	// var toc toc
	err = json.Unmarshal(body, v)
	// log.Print("ok")
	// log.Print(len(toc.Chapterlist.Chapters))
	// log.Print("o2k")

	return err
}
