package domain

type ArticleStates uint8

const (
	ArticleStatusUnknown ArticleStates = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (s ArticleStates) ToUint8() uint8 {
	return uint8(s)
}

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStates
}

func (a *Article) Abstract() string {
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return a.Content
	}
	return string(cs[:100])
}

type Author struct {
	Id   int64
	Name string
}
