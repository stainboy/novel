package zongheng

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"../util"
)

type Clawer interface {
	Process() error
}

func NewClawer() Clawer {
	return new(zhongheng)
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

	log.Println("Bye")
	return nil
}

func (c *zhongheng) initEnv() error {
	c.bookId = os.Getenv("bookId")
	c.bookTitle = os.Getenv("bookTitle")
	c.author = os.Getenv("author")
	c.pageSize, _ = strconv.Atoi(os.Getenv("pageSize"))
	c.secret = os.Getenv("bz")
	c.offset, _ = strconv.Atoi(os.Getenv("offset"))
	return nil
}

func (c *zhongheng) initToc() error {
	log.Printf("Fetching TOC...\n")
	url := fmt.Sprintf("http://m.zongheng.com/h5/ajax/chapter/list?h5=1&bookId=%s&pageNum=1&pageSize=%d&chapterId=0&asc=0", c.bookId, c.pageSize)
	// only cache for a day!
	key := fmt.Sprintf("toc_%s_%s", c.bookId, time.Now().UTC().Format("2006-01-02"))
	if err := util.FetchToc(url, key, &c.toc); err != nil {
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

	skipping := c.offset != 0
	for index, cp := range c.toc.Chapterlist.Chapters {
		if skipping {
			if cp.ChapterId != c.offset {
				continue
			} else {
				skipping = false
			}
		}

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
	skipping := c.offset != 0
	for _, cp := range c.toc.Chapterlist.Chapters {
		if skipping {
			if cp.ChapterId != c.offset {
				continue
			} else {
				skipping = false
			}
		}

		chapterName := fmt.Sprintf("##%s##\n\n", cp.ChapterName)
		util.AppendFile(bookMD, []byte(chapterName))

		c, err := ioutil.ReadFile(fmt.Sprintf("output/%s/tmp/%d.md", c.bookId, cp.ChapterId))
		if err != nil {
			return err
		}
		s := string(c)
		s = strings.Replace(s, "</p><p>（本章未完，请翻页）</p><p>", "", -1)
		s = strings.Replace(s, "<p>（本章完）</p>", "", -1)
		s = strings.Replace(s, "<p>", "", -1)
		s = strings.Replace(s, "</p>", "\n\n", -1)
		util.AppendFile(bookMD, []byte(s))
		util.AppendFile(bookMD, []byte("\n\n"))
	}

	return nil
}

func (c *zhongheng) produceMobi() error {
	log.Printf("Producing mobi file, it takes several minutes, please wait...")
	// https://manual.calibre-ebook.com/generated/en/ebook-convert.html
	shell := fmt.Sprintf("cd output/%s/;ebook-convert %s.md .mobi --title '%s' --rating 5 --authors '%s' --language zh --level1-toc '//h:h1' --level2-toc '//h:h2'", c.bookId, c.bookTitle, c.bookTitle, c.author)
	cmd := exec.Command("bash", "-c", shell)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
