package article

type Article struct {
	Id         int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title      string `gorm:"type=varchar(1024)" bson:"title,omitempty"`
	Content    string `gorm:"type=BLOB" bson:"content,omitempty"`
	AuthorId   int64  `gorm:"index=aid_ctime" bson:"authorId,omitempty"`
	Status     uint8  `bson:"status,omitempty"`
	CreateTime int64  `gorm:"index=aid_ctime" bson:"createTime,omitempty"`
	UpdateTime int64  `bson:"updateTime,omitempty"`
}

func (a *Article) TableName() string {
	return "article"
}

// PublishedArticle 发布文章
type PublishedArticle Article

func (a *PublishedArticle) TableName() string {
	return "publish_article"
}
