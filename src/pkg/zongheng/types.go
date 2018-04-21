package zongheng

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
	offset    int
	// ##
	toc toc
}
