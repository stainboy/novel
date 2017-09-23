package zongheng

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"../util"
)

type Clawer interface {
	Process() error
}

func NewClawer() Clawer {
	return new(zhongheng)
}

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

type chapter struct {
	Result struct {
		PageCount   int    `json:"pageCount"`
		ChapterNum  int    `json:"chapterNum"`
		ChapterName string `json:"chapterName"`
		ChapterId   string `json:"chapterId"`
		Content     string `json:"content"`
	} `json:"result"`
}

type zhongheng struct {
	// bookId=342974
	// bookTitle=永夜君王
	// author=烟雨江南
	// pageSize=2500
	// bz=342974|6122717|d8c8c2|aladin2_freexx
	bookId    string
	bookTitle string
	author    string
	pageSize  int
	secret    string
	// ##
	toc toc
}

func (c *zhongheng) Process() error {
	if err := c.initEnv(); err != nil {
		return err
	}

	if err := c.initToc(); err != nil {
		return err
	}

	if err := c.fetchChapters(); err != nil {
		return err
	}

	if err := c.produceMD(); err != nil {
		return err
	}

	if err := c.produceMobi(); err != nil {
		return err
	}

	return nil
}

func (c *zhongheng) initEnv() error {
	c.bookId = os.Getenv("bookId")
	c.bookTitle = os.Getenv("bookTitle")
	c.author = os.Getenv("author")
	c.pageSize, _ = strconv.Atoi(os.Getenv("pageSize"))
	c.secret = os.Getenv("bz")
	return nil
}

func (c *zhongheng) initToc() error {
	log.Printf("Fetching TOC...\n")
	url := fmt.Sprintf("http://m.zongheng.com/h5/ajax/chapter/list?h5=1&bookId=%s&pageNum=1&pageSize=%d&chapterId=0&asc=0", c.bookId, c.pageSize)
	if err := util.GetJSON(url, &c.toc); err != nil {
		return err
	}

	os.MkdirAll(fmt.Sprintf("output/%s/tmp", c.bookId), 0755)
	tocgo := fmt.Sprintf("output/%s/toc-go.json", c.bookId)
	r, err := json.Marshal(c.toc)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(tocgo, r, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *zhongheng) fetchChapters() error {

	cc := map[string]string{
		"___bz": c.secret,
	}

	total := len(c.toc.Chapterlist.Chapters)
	if os.Getenv("DEBUG_SIZE") != "" {
		total, _ = strconv.Atoi(os.Getenv("DEBUG_SIZE"))
	}
	for index, cp := range c.toc.Chapterlist.Chapters {

		// if cp.ChapterId != 4396038 {
		// 	continue
		// } else {
		// 	log.Printf("Debuging 4396038...")
		// }

		if index >= total {
			break
		}

		log.Printf("(%d/%d) Fetching chapter...\n", index+1, total)
		url := fmt.Sprintf("http://m.zongheng.com/h5/ajax/chapter?bookId=%s&chapterId=%d", c.bookId, cp.ChapterId)
		var head chapter
		err := util.GetJSONWithCookie(url, cc, &head)
		if err != nil {
			return err
		}

		chapterMD := fmt.Sprintf("output/%s/tmp/%d.md", c.bookId, cp.ChapterId)
		util.PurgeFile(chapterMD)

		count := head.Result.PageCount
		// log.Printf("Count -> %d\n", count)
		for i := 1; i <= count; i++ {
			if i == 1 {
				url = fmt.Sprintf("http://m.zongheng.com/h5/ajax/chapter?bookId=%s&chapterId=%d", c.bookId, cp.ChapterId)
			} else {
				url = fmt.Sprintf("http://m.zongheng.com/h5/ajax/chapter?bookId=%s&chapterId=%d_%d", c.bookId, cp.ChapterId, i)
			}
			var chapter chapter
			err := util.GetJSONWithCookie(url, cc, &chapter)
			if err != nil {
				return err
			}

			err = util.AppendFile(chapterMD, []byte(chapter.Result.Content))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *zhongheng) produceMD() error {

	log.Printf("Producing MD file, please wait...")

	bookMD := fmt.Sprintf("output/%s/%s.md", c.bookId, c.bookTitle)
	util.PurgeFile(bookMD)

	// title
	title := fmt.Sprintf("#%s#\n\n", c.bookTitle)
	util.AppendFile(bookMD, []byte(title))

	// chapters
	for _, cp := range c.toc.Chapterlist.Chapters {
		chapterName := fmt.Sprintf("##%s##\n\n", cp.ChapterName)
		util.AppendFile(bookMD, []byte(chapterName))

		c, err := ioutil.ReadFile(fmt.Sprintf("output/%s/tmp/%d.md", c.bookId, cp.ChapterId))
		if err != nil {
			return err
		}
		s := string(c)
		s = strings.Replace(s, "<p>", "", -1)
		s = strings.Replace(s, "</p>", "\n\n", -1)
		util.AppendFile(bookMD, []byte(s))
		util.AppendFile(bookMD, []byte("\n\n"))
	}

	return nil
}

func (c *zhongheng) produceMobi() error {
	log.Printf("(TODO) Producing mobi file, please wait...")
	// TODO
	return nil
}