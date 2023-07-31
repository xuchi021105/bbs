package model

import "gorm.io/gorm"

const (
	CommentType = "comment"
	ReplyType   = "reply"
)

type Reply struct {
	gorm.Model              // 回复ID和创建时间,以及实现软删除
	FromAvatarUrl    string // 用阿里云的oss存储的图片url(评论用户的头像)
	FromUserNickname string // 发出回复的用户昵称
	ToUserNickname   string // 被回复的用户昵称
	Content          string // 内容
	ReplyType        string // 用于判断是Comment还是Reply(是回复层主还是回复层中层) 这个用int类型性能会更好,但是懒,就算了
	CommentID        uint   // 发出回复关系的Comment外键,用于冗余,一次可以拿到所有在comment下的reply, gorm的外键
	FromUserID       uint   // 发出评论的UserID
	ToUserID         uint   // 被评论的UserID
}
