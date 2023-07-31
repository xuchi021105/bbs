package model

import (
	"backend/utils"
	"log"

	"gorm.io/gorm"
)

type Article struct {
	gorm.Model           // 文章ID和创建时间,以及实现软删除
	Title      string    // 文章标题
	UserID     uint      // 发布关系中的User外键
	Content    string    // 内容
	Toped      bool      // 由于判断是否置顶
	Majored    bool      // 用于判断是否是精华文章
	Comments   []Comment // 用户发布评论的关系
}

func (article *Article) IsExisted() bool {
	if err := utils.GetDB().First(&article).Error; err != nil {
		log.Printf("文章:%d不存在", article.ID)
		return false
	}
	return true
}
