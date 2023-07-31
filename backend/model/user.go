package model

import (
	"backend/service"
	"backend/utils"
	"fmt"
	"log"

	"gorm.io/gorm"
)

var (
	DefaultSignature            = "这个人很懒~,什么都没有写呢"
	DefaultUserAvatarPictureURL = fmt.Sprintf("%suserAvatarPicture/defaultUserAvatarPicture.png", service.OSSServerName)
)

type User struct {
	gorm.Model                 // 用户ID和用户的创建时间,以及实现软删除
	AvatarUrl        string    // 用阿里云的oss存储的图片url
	Nickname         string    // 昵称
	Signature        string    // 个性签名
	Role             string    // 角色 是管理员还是普通用户
	Password         string    // 密码
	Followers        []User    `gorm:"many2many:user_followers"` // 关注关系
	StaredArticles   []Article `gorm:"many2many:user_articles"`  // 收藏关系
	ReleasedArticles []Article // 发布关系
}

func (fromUser *User) IsFollowed(toUid uint) (bool, error) {
	var toUser User

	if err := utils.GetDB().Model(fromUser).Where("id = ?", toUid).Association("Followers").Find(&toUser); err != nil {
		log.Printf("数据库查询是否关注有问题: %+v", err)
		return false, err
	}
	if toUser.ID == toUid {
		log.Printf("用户%d关注用户%d", fromUser.ID, toUid)
		return true, nil
	} else {
		log.Printf("用户%d不关注用户%d", fromUser.ID, toUid)
		return false, nil
	}

}

// 下面这段自引用反向查询代码是从gorm的github issues里面找到的解决办法 https://github.com/bonfy/go-mega/issues/7
func (user *User) GetFansIDs() ([]uint, error) {
	var fansIDs []uint // 粉丝id的切片,就是关注我的人(FollowMePersonList)
	rows, err := utils.GetDB().Table("user_followers").Where("follower_id = ?", user.ID).Select("user_id, follower_id").Rows()
	if err != nil {
		// log.Println("Counting Following error:", err)
		log.Printf("数据库中查询user_followers表中粉丝ID有错误: %+v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, followerID int
		rows.Scan(&id, &followerID)
		fansIDs = append(fansIDs, uint(id))
	}
	return fansIDs, nil
}

func (user *User) IsExisted() bool {
	if err := utils.GetDB().First(&user).Error; err != nil {
		log.Printf("用户:%d未注册", user.ID)
		return false
	}
	return true
}

func (user *User) GetStaredArticleIDs() ([]uint, error) {

	var staredArticles []Article
	if err := utils.GetDB().Model(user).Association("StaredArticles").Find(&staredArticles); err != nil {
		log.Printf("数据库查询用户:%d的收藏文章出错", user.ID)
		return nil, err
	}

	staredArticleIDs := make([]uint, len(staredArticles))
	for _, staredAstaredArticle := range staredArticles {
		staredArticleIDs = append(staredArticleIDs, staredAstaredArticle.ID)
	}

	return staredArticleIDs, nil

}

func (user *User) IsStaredArticle(articleID uint) (bool, error) {
	var staredArticleIDs []uint
	var err error
	if staredArticleIDs, err = user.GetStaredArticleIDs(); err != nil {
		return false, err
	}
	staredArticleIDsSet := utils.ConvertUintSlice2Map(staredArticleIDs)
	_, ok := staredArticleIDsSet[articleID]
	return ok, nil
}
