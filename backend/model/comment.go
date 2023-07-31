package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model         // 评论ID和实现创建时间,以及实现软删除
	AvatarUrl  string  // 用阿里云的oss存储的图片url(评论用户的头像)
	Nickname   string  // 用户昵称
	Content    string  // 内容
	ArticleID  uint    // 评论发布关系的Article外键
	FromUserID uint    // 标明是那个用户发的评论
	Replys     []Reply // 发出回复关系
}
