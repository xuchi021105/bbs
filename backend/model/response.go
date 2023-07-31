package model

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type FollowerInfo struct {
	ID        uint   `json:"uid"`
	AvatarUrl string `json:"avatar_url"`
	Nickname  string `json:"nickname"`
	Signature string `json:"signature"`
}

type UserInfo struct {
	ID                    uint   `json:"uid"`
	AvatarUrl             string `json:"avatar_url"`
	Nickname              string `json:"nickname"`
	Signature             string `json:"signature"`
	CreateAtTime          string `json:"signInTime"`
	FollowersNumber       int    `json:"followers_number"`
	MyFollowersNumber     int    `json:"my_followers_number"` // fans数
	ArticlesNumber        int    `json:"articles_number"`
	StaredArticlesNumbser int    `json:"stared_articles_number"`
	IsMyself              bool   `json:"is_myself"`
	Followed              bool   `json:"followed"`
}

type ArticleInfo struct {
	IsMyArticle   bool   `json:"is_my_article"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Stared        bool   `json:"stared"`
	Toped         bool   `json:"toped"`
	Uid           uint   `json:"uid"`
	Nickname      string `json:"nickname"`
	Signature     string `json:"signature"`
	AvatarUrl     string `json:"avatar_url"`
	Followed      bool   `json:"followed"`
	ReleaseTime   string `json:"release_time"`
	CommentNumber uint   `json:"comment_number"`
	ArticleID     uint   `json:"article_id"`
}

type ArticleDigest struct {
	ArticleID   uint   `json:"article_id"`
	Title       string `json:"title"`
	ReleaseTime string `json:"release_time"`
	Toped       bool   `json:"toped"`
	Stared      bool   `json:"stared"`
	Majored     bool   `json:"majored"`
}

type ArticleDigestJSON struct {
	IsMyself bool            `json:"is_myself"`
	Articles []ArticleDigest `json:"articles"`
}

type StaredArticleDigest struct {
	Title       string `json:"title"`
	ArticleID   uint   `json:"article_id"`
	Uid         uint   `json:"uid"`
	Nickname    string `json:"nickname"`
	Majored     bool   `json:"majored"`
	ReleaseTime string `json:"release_time"`
	AvatarUrl   string `json:"avatar_url"`
}

type StaredArticleDigestJSON struct {
	IsMyself       bool                  `json:"is_myself"`
	StaredArticles []StaredArticleDigest `json:"stared_articles"`
}

type QueryArticleInfo struct {
	ArticleID uint   `json:"article_id"`
	Nickname  string `json:"nickname"`
	Uid       uint   `json:"uid"`
	AvatarUrl string `json:"avatar_url"`
	Title     string `json:"title"`
	Majored   bool   `json:"majored"`
}

type QueryUserInfo struct {
	AvatarUrl string `json:"avatar_url"`
	ID        uint   `json:"uid"`
	Nickname  string `json:"nickname"`
	Signature string `json:"signature"`
}

type CommentInfo struct {
	Uid         uint        `json:"uid"`
	CommentID   uint        `json:"comment_id"`
	Content     string      `json:"content"`
	AvatarUrl   string      `json:"avatar_url"`
	Nickname    string      `json:"nickname"`
	ReleaseTime string      `json:"release_time"`
	Replys      []ReplyInfo `json:"replys"`
}

type CommentResponseInfo struct {
	CommentID   uint   `json:"comment_id"`
	ReleaseTime string `json:"release_time"`
}

type ReplyInfo struct {
	ReplyType        string `json:"reply_type"`
	FromUid          uint   `json:"from_uid"`
	ToUid            uint   `json:"to_uid"`
	Content          string `json:"content"`
	FromAvatarUrl    string `json:"from_avatar_url"`
	ReplyID          uint   `json:"reply_id"`
	FromUserNickname string `json:"from_user_nickname"`
	ToUserNickname   string `json:"to_user_nickname"`
	ReleaseTime      string `json:"release_time"`
}

type ReplyResponseInfo struct {
	ReplyID     uint   `json:"reply_id"`
	ReleaseTime string `json:"release_time"`
}

var (
	SuccessCode          = 0 // TODO 这里应该用全大写来命名,之后再改吧
	DataQuestionCode     = 1
	UnAuthorizedCode     = 2
	ParamIsWrongCode     = 3
	ServerInnerErrorCode = 4

	UnAuthorizedResponse = Response{
		Code: UnAuthorizedCode,
		Msg:  "token verify failed",
	}

	ParamIsWrongResponse = Response{
		Code: ParamIsWrongCode,
		Msg:  "param is wrong",
	}

	ServerInnerErrorResponse = Response{
		Code: ServerInnerErrorCode,
		Msg:  "server inner error",
	}
)
