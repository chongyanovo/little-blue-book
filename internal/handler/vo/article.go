package vo

type ArticleVo struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Abstract   string `json:"abstract"`
	Content    string `json:"content"`
	AuthorId   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
	Status     uint8  `json:"status"`
}

type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type EditArticleRequest struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PublishArticleRequest struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ListArticleRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
