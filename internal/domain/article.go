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

type Author struct {
	Id int64
}
