package model

type Register struct {
	NickName string `json:"nickname"`
	Password string `json:"password"`
}

type Login struct {
	Uid      uint   `json:"uid"`
	Password string `json:"password"`
}

type Uid struct {
	Uid uint `json:"uid"`
}

type UpdateInfo struct {
	Nickname  string `json:"nickname"`
	Signature string `json:"signature"`
}

type ArticleID struct {
	ArticleID uint `json:"article_id"`
}

type ArticleJSON struct {
	ArticleID uint   `json:"article_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Toped     bool   `json:"toped"`
}

type CommentJSON struct {
	ArticleID uint   `json:"article_id"`
	Content   string `json:"content"`
}

type ReplyJSON struct {
	CommentID uint   `json:"comment_id"`
	Content   string `json:"content"`
	ReplyType string `json:"reply_type"`
	ToUid     uint   `json:"to_uid"`
}

type KeywordJSON struct {
	Keyword string `json:"keyword"`
}
