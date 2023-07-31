package handler

//TODO 记得给数据加上分页功能(利用page和pageSize,得到offset,然后利用这个得到分页的数据)

import (
	"backend/config"
	"backend/model"
	"backend/service"
	"backend/utils"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Login(context *gin.Context) {

	var loginJson model.Login

	if err := context.ShouldBindJSON(&loginJson); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}

	db := utils.GetDB()

	var user model.User

	if err := db.First(&user, loginJson.Uid).Error; err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code: model.DataQuestionCode,
			Msg:  "user id does not exist",
		})
		log.Printf("用户id不存在,未注册: %+v", err)
		return
	}

	if user.Password != utils.GetSHA256HashCode(loginJson.Password) {
		context.JSON(http.StatusOK, model.Response{
			Code: model.DataQuestionCode,
			Msg:  "password is wrong",
		})
		log.Printf("用户密码错误")
	} else {
		if tokenString, err := utils.GeneToken(loginJson.Uid); err != nil {
			context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
			log.Printf("签发token失败: %+v", err)
			return
		} else {
			context.JSON(http.StatusOK, model.Response{
				Code: model.SuccessCode,
				Msg:  "success",
				Data: gin.H{
					"token":     tokenString,
					"uid":       user.ID,
					"nickname":  user.Nickname,
					"avatarUrl": user.AvatarUrl,
					"signature": user.Signature,
				},
			})
			log.Printf("用户登录成功, 签发token,token为: %s", tokenString)
		}
	}

}

func Register(context *gin.Context) {
	var registerJson model.Register

	if err := context.ShouldBindJSON(&registerJson); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}

	db := utils.GetDB()
	user := model.User{
		Nickname:  registerJson.NickName,
		Password:  utils.GetSHA256HashCode(registerJson.Password),
		Signature: model.DefaultSignature,
		Role:      `User`,
	}

	if err := db.Create(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("注册时数据库创建用户失败: %+v", err)
		return
	}

	user.AvatarUrl = fmt.Sprintf("%suserAvatarPicture/%d", service.OSSServerName, user.ID)

	if err := db.Save(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("保存时数据库操作失败: %+v", err)
		return
	}

	if err := service.GetAvatarPictureBucket().PutObjectFromFile(fmt.Sprintf("userAvatarPicture/%d", user.ID), fmt.Sprintf("%s/resources/images/defaultUserAvatarPicture.png", config.RootDir)); err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("OSS上传文件失败")
		return
	}

	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
		Data: gin.H{
			"uid": user.ID,
		},
	})
	log.Printf("用户注册成功,昵称为:%s, 用户ID为:%d", registerJson.NickName, user.ID)
}

func GetFollowers(context *gin.Context) {

	var uidJson model.Uid
	if err := context.ShouldBindJSON(&uidJson); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}

	user := &model.User{}
	user.ID = uidJson.Uid
	var followerInfos []model.FollowerInfo // 控制查询出的字段,也可以通过Select函数来控制,这里是用结构体智能控制
	if user.IsExisted() {
		// 注释掉的这段是用Preload来,两个都试试吧,看看哪个是可行的
		// if err := utils.GetDB().Model(&model.User{}).Where("id = ?", uidJson.Uid).Preload("Users").Find(&followerInfos).Error; err != nil {
		if err := utils.GetDB().Model(&user).Association("Followers").Find(&followerInfos); err != nil {
			context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
			log.Printf("数据库出错,未找到用户id: %d 的关注者: %+v", uidJson.Uid, err)
			return
		}
		if len(followerInfos) == 0 {
			context.JSON(http.StatusOK, model.Response{
				Code: model.SuccessCode,
				Msg:  "success but no data",
			})
			log.Printf("关注者为0")
		} else {
			context.JSON(http.StatusOK, model.Response{
				Code: model.SuccessCode,
				Msg:  "success",
				Data: followerInfos,
			})
			log.Printf("成功返回关注者信息")
		}
	} else {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("用户不存在")
		return
	}
}

func GetFollowMePersonList(context *gin.Context) {

	var uidJson model.Uid

	if err := context.ShouldBindJSON(&uidJson); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}

	user := model.User{}
	user.ID = uidJson.Uid
	if user.IsExisted() {
		var fansIDs []uint // 粉丝id的切片,就是关注我的人(FollowMePersonList)
		fansIDs, err := user.GetFansIDs()
		if err != nil {
			// log.Println("Counting Following error:", err)
			context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
			log.Printf("数据库中查询user_followers表中粉丝ID有错误: %+v", err)
			return
		}

		log.Print("fans:", fansIDs)

		var fansInfos []model.FollowerInfo

		if len(fansIDs) == 0 {

			context.JSON(http.StatusOK, model.Response{
				Code: model.SuccessCode,
				Msg:  "success but no data",
			})
			log.Print("粉丝数为0")

		} else {

			if err := utils.GetDB().Model(&model.User{}).Find(&fansInfos, fansIDs).Error; err != nil {
				context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
				log.Printf("数据库通过粉丝ID查询信息出错: %+v", err)
				return
			}
			context.JSON(http.StatusOK, model.Response{
				Code: model.SuccessCode,
				Msg:  "success",
				Data: fansInfos,
			})
			log.Printf("成功返回粉丝信息")
		}
	} else {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("用户不存在")
		return
	}
}

func FollowUser(context *gin.Context) {
	var uidJson model.Uid
	if err := context.ShouldBindJSON(&uidJson); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}
	fromUserIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	fromUserID, ok := fromUserIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}
	fromUser := model.User{}
	fromUser.ID = fromUserID
	toUser := model.User{}
	toUser.ID = uidJson.Uid

	// 这里用toUser.IsExisted()函数就好
	if err := utils.GetDB().First(&toUser).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("uid为:%d的用户在数据库中不存在", uidJson.Uid)
		return
	}

	if uidJson.Uid == fromUserID {
		context.JSON(http.StatusOK, model.Response{
			Code: model.DataQuestionCode,
			Msg:  "can't follow myself",
		})
		log.Print("用户不能自己关注自己")
		return
	}
	if isFollowed, err := fromUser.IsFollowed(uidJson.Uid); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("数据库查询关注时有问题")
	} else {
		if !isFollowed { // 未关注
			if err := utils.GetDB().Model(&fromUser).Association("Followers").Append(&toUser); err != nil {
				context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
				log.Print("给用户通过Association方法来增加关联失败")
				return
			} else {
				context.JSON(http.StatusOK, model.Response{
					Code: model.SuccessCode,
					Msg:  "success",
				})
				log.Printf("用户:%d(fromUser)成功关注用户:%d(toUser)", fromUserID, toUser.ID)
				return
			}
		} else { // 已经关注
			context.JSON(http.StatusOK, model.Response{
				Code: model.DataQuestionCode,
				Msg:  "already follow",
			})
			log.Printf("用户:%d(fromUser)已经关注用户:%d(toUser),不能重复关注", fromUserID, toUser.ID)
			return
		}
	}
}

func CancelFollowedUser(context *gin.Context) {
	var uidJson model.Uid
	if err := context.ShouldBindJSON(&uidJson); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}
	fromUserIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	fromUserID, ok := fromUserIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}
	fromUser := model.User{}
	fromUser.ID = fromUserID
	toUser := model.User{}
	toUser.ID = uidJson.Uid
	if err := utils.GetDB().Find(&toUser).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("uid为:%d的用户在数据库中不存在", uidJson.Uid)
		return
	}
	if isFollowed, err := fromUser.IsFollowed(uidJson.Uid); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("数据库查询关注时有问题")
	} else {
		if !isFollowed { // 未关注
			context.JSON(http.StatusOK, model.Response{
				Code: model.DataQuestionCode,
				Msg:  "no follow before",
			})
			log.Printf("用户:%d(fromUser)未关注用户:%d(toUser),不能取消没有关注的关注", fromUserID, toUser.ID)
			return
		} else { // 已经关注
			if err := utils.GetDB().Model(&fromUser).Association("Followers").Delete(&toUser); err != nil {
				context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
				log.Print("给用户通过Association方法来增加关联失败")
				return
			} else {
				context.JSON(http.StatusOK, model.Response{
					Code: model.SuccessCode,
					Msg:  "success",
				})
				log.Printf("用户:%d(fromUser)成功取消关注用户:%d(toUser)", fromUserID, toUser.ID)
				return
			}
		}
	}
}

func GetUserInfoByUid(context *gin.Context) {
	var uidJson model.Uid
	if err := context.ShouldBindJSON(&uidJson); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}

	var fromUserID uint
	fromUserIDAny, ok := context.Get("myUid")
	tokenIsExisted := false
	if ok { // 有token
		log.Print("找到token中的myUid")
		fromUserID, ok = fromUserIDAny.(uint)
		if !ok { // 失败
			log.Print("传来的myUid不是uint类型,类型断言失败")
		} else {
			tokenIsExisted = true
		}
	} else { // 没有token的情况
		log.Print("没有找到token中的myUid")
	}

	user := model.User{}
	user.ID = uidJson.Uid
	if err := utils.GetDB().Preload("Followers").Preload("StaredArticles").Preload("ReleasedArticles").First(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("数据库中不存在用户:%d", user.ID)
		return
	}

	var fansIDs []uint // 粉丝id的切片,就是关注我的人(FollowMePersonList)
	fansIDs, err := user.GetFansIDs()
	if err != nil {
		// log.Println("Counting Following error:", err)
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("数据库中查询user_followers表中粉丝ID有错误: %+v", err)
		return
	}

	userInfo := model.UserInfo{
		ID:                    user.ID,
		AvatarUrl:             user.AvatarUrl,
		Nickname:              user.Nickname,
		Signature:             user.Signature,
		CreateAtTime:          user.CreatedAt.String(),
		FollowersNumber:       len(user.Followers),
		MyFollowersNumber:     len(fansIDs),
		ArticlesNumber:        len(user.ReleasedArticles),
		StaredArticlesNumbser: len(user.StaredArticles),
	}
	if tokenIsExisted {

		if fromUserID == uidJson.Uid {
			userInfo.IsMyself = true
		} else {
			userInfo.IsMyself = false
		}

		for _, fansID := range fansIDs {
			if fansID == fromUserID {
				userInfo.Followed = true
				break
			}
		}

	} else {
		userInfo.IsMyself = false
	}

	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
		Data: userInfo,
	})
	log.Printf("成功查询到用户%d的信息", userInfo.ID)

}

// 更新的时候只需要判断提交的字段就好(判断是否是空),但是更新图片的时候需要文件处理,放在另一个函数里面了
func UpdateUserInfo(context *gin.Context) {

	var updateInfo model.UpdateInfo
	if err := context.ShouldBindJSON(&updateInfo); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}

	userIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}

	user := model.User{}
	user.ID = userID

	if err := utils.GetDB().First(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("uid为:%d的用户在数据库中不存在", user.ID)
		return
	}

	if updateInfo.Nickname != "" {
		user.Nickname = updateInfo.Nickname
	}
	if updateInfo.Signature != "" {
		user.Signature = updateInfo.Signature
	}

	if err := utils.GetDB().Save(user).Error; err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("数据库储存用户:%d失败", user.ID)
		return
	}

	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
	})
	log.Printf("用户:%d更新数据成功", user.ID)

}

// TODO 这里还应该考虑文件不是图片文件类型的情况,但是懒,有时间再做吧
func UploadAvatarPicture(context *gin.Context) {

	userIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}

	user := model.User{}
	user.ID = userID

	if err := utils.GetDB().First(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("uid为:%d的用户在数据库中不存在", user.ID)
		return
	}

	var fileHeader *multipart.FileHeader
	var err error

	if fileHeader, err = context.FormFile("file"); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("请求参数有误, 没有解析出文件")
		return
	}
	var file multipart.File

	// Open函数会创建一个临时文件并且打开(os.Open)
	if file, err = fileHeader.Open(); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("打开文件失败")
		return
	}

	defer file.Close()

	avatarUrl := fmt.Sprintf("userAvatarPicture/%d", user.ID)

	if err := service.GetAvatarPictureBucket().PutObject(avatarUrl, file); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("上传图片至oss失败")
		return
	}

	user.AvatarUrl = service.OSSServerName + avatarUrl
	if err := utils.GetDB().Save(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("avatarUrl保存失败")
		return
	}

	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  `success`,
	})
	log.Printf("用户:%d更新头像成功", user.ID)

}

func PostArticle(context *gin.Context) {
	var articleJSON model.ArticleJSON
	if err := context.ShouldBindJSON(&articleJSON); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("json解析有误,参数错误")
		return
	}

	userIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}

	article := model.Article{}
	article.Title = articleJSON.Title
	article.UserID = userID
	article.Content = articleJSON.Content
	article.Toped = articleJSON.Toped

	if err := utils.GetDB().Create(&article).Error; err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("用户:%d在发布文章时数据库出错", userID)
		return
	}
	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
	})
	log.Printf("用户:%d发布文章:%d成功", userID, article.ID)
}

func UpdateArticle(context *gin.Context) {

	var articleJSON model.ArticleJSON
	if err := context.ShouldBindJSON(&articleJSON); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("json解析有误,参数错误")
		return
	}

	userIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}

	article := model.Article{}
	article.ID = articleJSON.ArticleID

	if !article.IsExisted() {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("文章:%d不存在,不能更新", article.ID)
		return
	}

	if userID != article.UserID {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("要更新的文章:%d不属于用户:%d所有,属于用户:%d,不能更新", article.ID, userID, article.UserID)
		return
	}

	article.Title = articleJSON.Title
	article.UserID = userID
	article.Content = articleJSON.Content
	article.Toped = articleJSON.Toped

	if err := utils.GetDB().Save(&article).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongCode)
		log.Printf("数据库更新文章:%d时有误", article.ID)
		return
	}

	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
	})
	log.Printf("用户:%d成功更新文章:%d", article.UserID, article.ID)
}

func DeleteArticle(context *gin.Context) {

	var articleID model.ArticleID
	if err := context.ShouldBindJSON(&articleID); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("json解析有误,参数错误")
		return
	}

	userIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}
	article := model.Article{}
	article.ID = articleID.ArticleID

	if !article.IsExisted() {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("文章:%d不存在,不能删除", article.ID)
		return
	}

	if userID != article.UserID {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("要删除的文章:%d不属于用户:%d所有,属于用户:%d,不能删除", article.ID, userID, article.UserID)
		return
	}

	if err := utils.GetDB().Delete(&article).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("数据库删除文章:%d时有误", article.ID)
		return
	}

	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
	})
	log.Printf("用户:%d成功删除文章:%d", article.UserID, article.ID)
}

func GetArticleInfo(context *gin.Context) {
	var articleID model.ArticleID
	if err := context.ShouldBindJSON(&articleID); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("json解析有误,参数错误")
		return
	}

	article := model.Article{}
	article.ID = articleID.ArticleID

	if err := utils.GetDB().Preload("Comments").First(&article).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("数据库中不存在文章:%d", article.ID)
		return
	}

	var isMyArticle bool
	stared := false
	var userID uint

	var err error

	userIDAny, ok := context.Get("myUid")
	if !ok { // 没有token的情况
		log.Print("没有找到token中的myUid")
		isMyArticle = false
	} else { // 有token
		userID, ok = userIDAny.(uint)
		if !ok { // 失败
			log.Print("传来的myUid不是uint类型,类型断言失败")
			isMyArticle = false
		} else {
			user := model.User{}
			user.ID = userID
			if stared, err = user.IsStaredArticle(articleID.ArticleID); err != nil {
				context.JSON(http.StatusOK, model.ParamIsWrongResponse)
				log.Printf("查询用户:%d的收藏文章时出错", user.ID)
				return
			}

			if userID == article.UserID {
				isMyArticle = true
			} else {
				isMyArticle = false
			}
		}
	}

	user := model.User{}     // 作者的信息
	user.ID = article.UserID // 作者的id
	if err := utils.GetDB().First(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("用户:%d不存在", userID)
		return
	}

	followed := false

	if isMyArticle { // 是自己的文章
		followed = false
	} else { // 不是自己的文章
		var fansIDs []uint
		var err error
		if fansIDs, err = user.GetFansIDs(); err != nil {
			context.JSON(http.StatusOK, model.ParamIsWrongResponse)
			log.Printf("查找用户:%d的关注列表失败", userID)
			return
		}
		for _, fansID := range fansIDs {
			if fansID == userID {
				followed = true
				break
			}
		}
	}

	articleInfo := model.ArticleInfo{
		IsMyArticle:   isMyArticle,
		Title:         article.Title,
		Content:       article.Content,
		Stared:        stared,
		Toped:         article.Toped,
		Uid:           article.UserID,
		Nickname:      user.Nickname,
		Signature:     user.Signature,
		AvatarUrl:     user.AvatarUrl,
		Followed:      followed,
		ReleaseTime:   article.CreatedAt.String(),
		CommentNumber: uint(len(article.Comments)),
		ArticleID:     article.ID,
	}

	context.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "success",
		Data: articleInfo,
	})
	log.Printf("成功查询到文章%d的信息", articleInfo.ArticleID)
}

func GetUserAllArticles(context *gin.Context) {

	var uidJson model.Uid
	if err := context.ShouldBindJSON(&uidJson); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}

	takeArticleuser := model.User{} // 要获取文章的用户
	takeArticleuser.ID = uidJson.Uid

	if isExisted := takeArticleuser.IsExisted(); !isExisted {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("用户:%d不存在", takeArticleuser.ID)
		return
	}

	var articles []model.Article
	if err := utils.GetDB().Model(&takeArticleuser).Association("ReleasedArticles").Find(&articles); err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("数据库查询出错,未找到用户:%d下的文章", uidJson.Uid)
		return
	}

	isMyself := false
	var userID uint
	var user model.User // 带token的自己
	userIDAny, ok := context.Get("myUid")
	if !ok { // 没有token的情况
		log.Print("没有找到token中的myUid")
		isMyself = false
	} else { // 有token
		userID, ok = userIDAny.(uint)
		if !ok { // 失败
			log.Print("传来的myUid不是uint类型,类型断言失败")
			isMyself = false
		} else { // 成功
			user.ID = userID
			if userID != uidJson.Uid {
				isMyself = false
			} else {
				isMyself = true
			}
		}
	}
	var staredArticleIDs []uint
	var staredArticleIDsSet map[uint]struct{}
	var err error
	if ok {
		if staredArticleIDs, err = user.GetStaredArticleIDs(); err != nil {
			context.JSON(http.StatusOK, model.ParamIsWrongResponse)
			log.Printf("查询用户:%d的收藏文章时出现问题", user.ID)
			return
		}
		staredArticleIDsSet = utils.ConvertUintSlice2Map(staredArticleIDs)
	}

	var articleDigests []model.ArticleDigest
	for _, article := range articles {
		var stared bool
		if ok {
			_, stared = staredArticleIDsSet[article.ID]
		}
		articleDigest := model.ArticleDigest{
			ArticleID:   article.ID,
			Title:       article.Title,
			ReleaseTime: article.CreatedAt.String(),
			Toped:       article.Toped,
			Stared:      stared,
			Majored:     article.Majored,
		}
		articleDigests = append(articleDigests, articleDigest)
	}

	if len(articleDigests) == 0 {
		context.JSON(http.StatusOK, model.Response{
			Code: model.DataQuestionCode,
			Msg:  "success but no data",
		})
		log.Printf("用户:%d下的文章为0", uidJson.Uid)
	} else {
		context.JSON(http.StatusOK, model.Response{
			Code: model.SuccessCode,
			Msg:  "success",
			Data: model.ArticleDigestJSON{
				IsMyself: isMyself,
				Articles: articleDigests,
			},
		})
		log.Printf("成功查询用户:%d下的文章", uidJson.Uid)
	}

}

func StarArticle(context *gin.Context) {

	var articleID model.ArticleID
	if err := context.ShouldBindJSON(&articleID); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("json解析有误,参数错误")
		return
	}

	UserIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	UserID, ok := UserIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}

	user := model.User{}
	user.ID = UserID

	article := model.Article{}
	article.ID = articleID.ArticleID

	// 看要收藏的文章是否存在
	if err := utils.GetDB().First(&article).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("文章:%d在数据库中不存在", article.ID)
		return
	}

	if stared, err := user.IsStaredArticle(article.ID); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("数据库查询文章时有问题")
	} else {
		if !stared { // 未收藏
			if err := utils.GetDB().Model(&user).Association("StaredArticles").Append(&article); err != nil {
				context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
				log.Printf("给用户:%d通过Association方法来收藏文章%d失败", user.ID, article.ID)
				return
			} else {
				context.JSON(http.StatusOK, model.Response{
					Code: model.SuccessCode,
					Msg:  "success",
				})
				log.Printf("用户:%d成功收藏文章:%d", user.ID, article.ID)
				return
			}
		} else { // 已经收藏
			context.JSON(http.StatusOK, model.Response{
				Code: model.DataQuestionCode,
				Msg:  "already star",
			})
			log.Printf("用户:%d已经收藏文章:%d,不能重复收藏", user.ID, article.ID)
			return
		}
	}
}

func CancelStaredArticle(context *gin.Context) {
	var articleID model.ArticleID
	if err := context.ShouldBindJSON(&articleID); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("json解析有误,参数错误")
		return
	}

	UserIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	UserID, ok := UserIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}

	user := model.User{}
	user.ID = UserID

	article := model.Article{}
	article.ID = articleID.ArticleID

	// 看要取消收藏的文章是否存在
	if err := utils.GetDB().First(&article).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("文章:%d在数据库中不存在", article.ID)
		return
	}

	if stared, err := user.IsStaredArticle(article.ID); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("数据库查询文章时有问题")
	} else {
		if !stared { // 未收藏
			context.JSON(http.StatusOK, model.Response{
				Code: model.DataQuestionCode,
				Msg:  "don't star",
			})
			log.Printf("用户:%d没有收藏文章:%d,不能取消收藏", user.ID, article.ID)
			return
		} else { // 已经收藏
			if err := utils.GetDB().Model(&user).Association("StaredArticles").Delete(&article); err != nil {
				context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
				log.Printf("给用户:%d通过Association方法来取消收藏文章%d失败", user.ID, article.ID)
				return
			} else {
				context.JSON(http.StatusOK, model.Response{
					Code: model.SuccessCode,
					Msg:  "success",
				})
				log.Printf("用户:%d成功取消收藏文章:%d", user.ID, article.ID)
				return
			}
		}
	}
}

func GetUserAllStaredArticles(context *gin.Context) {

	var uidJson model.Uid
	if err := context.ShouldBindJSON(&uidJson); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误,解析json失败")
		return
	}

	takeStaredArticleuser := model.User{} // 要获取文章的用户
	takeStaredArticleuser.ID = uidJson.Uid

	if err := utils.GetDB().First(&takeStaredArticleuser).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("用户:%d不存在", takeStaredArticleuser.ID)
		return
	}

	var articles []model.Article
	if err := utils.GetDB().Model(&takeStaredArticleuser).Association("StaredArticles").Find(&articles); err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("数据库查询出错,未找到用户:%d下收藏的文章", uidJson.Uid)
		return
	}

	isMyself := false
	var userID uint
	var user model.User // 带token的自己
	userIDAny, ok := context.Get("myUid")
	if !ok { // 没有token的情况
		log.Print("没有找到token中的myUid")
		isMyself = false
	} else { // 有token
		userID, ok = userIDAny.(uint)
		if !ok { // 失败
			log.Print("传来的myUid不是uint类型,类型断言失败")
			isMyself = false
		} else { // 成功
			user.ID = userID
			if userID != uidJson.Uid {
				isMyself = false
			} else {
				isMyself = true
			}
		}
	}

	var staredArticleDigests []model.StaredArticleDigest
	for _, article := range articles {
		staredArticleDigest := model.StaredArticleDigest{
			Title:       article.Title,
			ArticleID:   article.ID,
			Uid:         article.UserID,
			Nickname:    takeStaredArticleuser.Nickname,
			Majored:     article.Majored,
			ReleaseTime: article.CreatedAt.String(),
			AvatarUrl:   takeStaredArticleuser.AvatarUrl,
		}
		staredArticleDigests = append(staredArticleDigests, staredArticleDigest)
	}

	if len(staredArticleDigests) == 0 {
		context.JSON(http.StatusOK, model.Response{
			Code: model.DataQuestionCode,
			Msg:  "success but no data",
		})
		log.Printf("用户:%d下的收藏文章为0", uidJson.Uid)
	} else {
		context.JSON(http.StatusOK, model.Response{
			Code: model.SuccessCode,
			Msg:  "success",
			Data: model.StaredArticleDigestJSON{
				IsMyself:       isMyself,
				StaredArticles: staredArticleDigests,
			},
		})
		log.Printf("成功查询用户:%d下收藏的文章", uidJson.Uid)
	}

}

func CommentArticle(context *gin.Context) {
	var commentJSON model.CommentJSON
	if err := context.ShouldBindJSON(&commentJSON); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误")
		return
	}

	userIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}

	user := model.User{}
	user.ID = userID
	if err := utils.GetDB().First(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("用户:%d不存在", user.ID)
		return
	}

	article := model.Article{}
	article.ID = commentJSON.ArticleID
	if isExisted := article.IsExisted(); !isExisted {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("文章:%d不存在,不能评论", article.ID)
		return
	}

	comment := model.Comment{}
	comment.AvatarUrl = user.AvatarUrl
	comment.Nickname = user.Nickname
	comment.Content = commentJSON.Content
	comment.ArticleID = article.ID
	comment.FromUserID = user.ID

	if err := utils.GetDB().Create(&comment).Error; err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("创建评论时失败")
		return
	}

	if err := utils.GetDB().Model(&article).Association("Comments").Append(&comment); err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("用户:%d给文章:%d,添加评论:%d失败", user.ID, article.ID, comment.ID)
		return
	}
	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
		Data: model.CommentResponseInfo{
			CommentID:   comment.ID,
			ReleaseTime: comment.CreatedAt.String(),
		},
	})
	log.Printf("用户:%d给文章:%d,添加评论:%d成功", user.ID, article.ID, comment.ID)
}

// 用@来区分Reply中的回复
// TODO 后续可能会因为评论和回复太多,所以要进行lazy load,分页进行加载,一次加载10条这种,但是这是后面的事情
func GetAllCommentsAndReplysInArticle(context *gin.Context) {

	var ArticleID model.ArticleID
	if err := context.ShouldBindJSON(&ArticleID); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误")
		return
	}
	article := model.Article{}
	article.ID = ArticleID.ArticleID
	if isExisted := article.IsExisted(); !isExisted {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("文章%d不存在", ArticleID.ArticleID)
		return
	}

	var comments []model.Comment
	if err := utils.GetDB().Model(&article).Association("Comments").Find(&comments); err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("查询文章:%d下的评论失败", article.ID)
		return
	}

	if len(comments) == 0 {
		context.JSON(http.StatusOK, model.Response{
			Code: model.SuccessCode,
			Msg:  "success but no data",
		})
		log.Printf("文章:%d下没有评论", ArticleID.ArticleID)
		return
	}
	commentInfos := make([]model.CommentInfo, 0)

	for _, comment := range comments {
		commentInfo := model.CommentInfo{}
		commentInfo.Uid = comment.FromUserID
		commentInfo.CommentID = comment.ID
		commentInfo.Content = comment.Content
		commentInfo.AvatarUrl = comment.AvatarUrl
		commentInfo.Nickname = comment.Nickname
		commentInfo.ReleaseTime = comment.CreatedAt.String()

		var replys []model.Reply
		if err := utils.GetDB().Model(&comment).Association("Replys").Find(&replys); err != nil {
			context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
			log.Printf("查询文章:%d下的评论:%d下的回复失败", article.ID, comment.ID)
			return
		}
		replyInfos := make([]model.ReplyInfo, 0)
		for _, reply := range replys {
			replyInfo := model.ReplyInfo{}
			replyInfo.ReplyType = reply.ReplyType
			replyInfo.FromUid = reply.FromUserID
			replyInfo.Content = reply.Content
			replyInfo.FromAvatarUrl = reply.FromAvatarUrl
			replyInfo.ReplyID = reply.ID
			replyInfo.FromUserNickname = reply.FromUserNickname
			replyInfo.ReleaseTime = reply.CreatedAt.String()
			if replyInfo.ReplyType == model.CommentType { // 针对评论的回复(无@的人)

			} else {
				if replyInfo.ReplyType == model.ReplyType { // 针对回复的回复(有@的人)
					replyInfo.ToUid = reply.ToUserID
					replyInfo.ToUserNickname = reply.ToUserNickname
				}
			}
			replyInfos = append(replyInfos, replyInfo)
		}
		commentInfo.Replys = replyInfos
		commentInfos = append(commentInfos, commentInfo)
	}

	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
		Data: commentInfos,
	})
	log.Printf("获取文章:%d下的所有评论回复成功", article.ID)

}

func ReplyComment(context *gin.Context) {
	var replyJSON model.ReplyJSON
	if err := context.ShouldBindJSON(&replyJSON); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("参数有误")
		return
	}

	userIDAny, ok := context.Get("myUid")
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("没有找到token中的myUid")
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Print("传来的myUid不是uint类型,类型断言失败")
		return
	}

	user := model.User{}
	user.ID = userID
	if err := utils.GetDB().First(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("用户:%d不存在", user.ID)
		return
	}

	comment := model.Comment{}
	comment.ID = replyJSON.CommentID
	if err := utils.GetDB().First(&comment).Error; err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("评论:%d不存在,回复失败", comment.ID)
		return
	}

	reply := model.Reply{}
	reply.CommentID = replyJSON.CommentID
	reply.FromAvatarUrl = user.AvatarUrl
	reply.FromUserNickname = user.Nickname
	reply.Content = replyJSON.Content
	reply.FromUserID = user.ID

	if replyJSON.ReplyType == model.CommentType { // 针对评论

		reply.ReplyType = model.CommentType

	} else {
		if replyJSON.ReplyType == model.ReplyType { // 针对回复(有要回复@的人)

			reply.ReplyType = model.ReplyType
			toUser := model.User{}
			toUser.ID = replyJSON.ToUid
			if err := utils.GetDB().First(&toUser).Error; err != nil {
				context.JSON(http.StatusOK, model.ParamIsWrongResponse)
				log.Printf("用户:%d不存在", toUser.ID)
				return
			}
			reply.ToUserNickname = toUser.Nickname
			reply.ToUserID = toUser.ID

		} else {
			context.JSON(http.StatusOK, model.ParamIsWrongResponse)
			log.Printf("reply_type类型有错")
			return
		}
	}

	if err := utils.GetDB().Create(&reply).Error; err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("创建回复时失败")
		return
	}

	if err := utils.GetDB().Model(&comment).Association("Replys").Append(&reply); err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("给评论:%d加上回复:%d时失败", comment.ID, reply.ID)
		return
	}

	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
		Data: model.ReplyResponseInfo{
			ReplyID:     reply.ID,
			ReleaseTime: reply.CreatedAt.String(),
		},
	})
	log.Printf("用户:%d给评论:%d,添加回复:%d成功", user.ID, comment.ID, comment.ID)

}

// 模糊查询,只能单关键字(用like来做)
// TODO 多关键字以后有机会来做吧
func QueryArticleByKeyword(context *gin.Context) {
	var KeywordJSON model.KeywordJSON
	if err := context.ShouldBindJSON(&KeywordJSON); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("JSON解析有误")
		return
	}

	var articles []model.Article

	if err := utils.GetDB().Where("title like ?", "%"+KeywordJSON.Keyword+"%").Find(&articles).Error; err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("模糊查询文章失败")
		return
	}

	if len(articles) == 0 {
		context.JSON(http.StatusOK, model.Response{
			Code: model.DataQuestionCode,
			Msg:  "no data",
		})
		log.Printf("没有根据关键词:%s查询到文章", KeywordJSON.Keyword)
		return
	}
	queryArticleInfos := make([]model.QueryArticleInfo, len(articles))
	for index, article := range articles {
		user := model.User{}
		user.ID = article.UserID
		if err := utils.GetDB().First(&user).Error; err != nil {
			context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
			log.Printf("用户:%d不存在,查询文章失败", user.ID)
			return
		}
		queryAricleInfo := model.QueryArticleInfo{
			ArticleID: article.ID,
			Nickname:  user.Nickname,
			Uid:       user.ID,
			AvatarUrl: user.AvatarUrl,
			Title:     article.Title,
			Majored:   article.Majored,
		}
		queryArticleInfos[index] = queryAricleInfo
	}
	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
		Data: queryArticleInfos,
	})
	log.Printf("成功根据关键字:%s查询到文章列表", KeywordJSON.Keyword)
}

// 模糊查询
func QueryUserByKeyword(context *gin.Context) {

	var KeywordJSON model.KeywordJSON
	if err := context.ShouldBindJSON(&KeywordJSON); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("JSON解析有误")
		return
	}

	var users []model.User

	if err := utils.GetDB().Where("nickname like ?", "%"+KeywordJSON.Keyword+"%").Find(&users).Error; err != nil {
		context.JSON(http.StatusOK, model.ServerInnerErrorResponse)
		log.Printf("模糊查询用户昵称失败")
		return
	}

	if len(users) == 0 {
		context.JSON(http.StatusOK, model.Response{
			Code: model.DataQuestionCode,
			Msg:  "no data",
		})
		log.Printf("没有根据关键词:%s查询到用户", KeywordJSON.Keyword)
		return
	}
	queryUserInfos := make([]model.QueryUserInfo, 0)
	for _, user := range users {
		queryUserInfo := model.QueryUserInfo{
			AvatarUrl: user.AvatarUrl,
			ID:        user.ID,
			Nickname:  user.Nickname,
			Signature: user.Signature,
		}
		queryUserInfos = append(queryUserInfos, queryUserInfo)
	}
	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
		Data: queryUserInfos,
	})
	log.Printf("成功根据关键字:%s查询到用户列表", KeywordJSON.Keyword)
}

// 精确查询
func QueryUserByUID(context *gin.Context) {

	var KeywordJSON model.KeywordJSON
	if err := context.ShouldBindJSON(&KeywordJSON); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("JSON解析有误")
		return
	}

	user := model.User{}
	var uid int
	var err error
	if uid, err = strconv.Atoi(KeywordJSON.Keyword); err != nil {
		context.JSON(http.StatusOK, model.ParamIsWrongResponse)
		log.Printf("参数格式有误")
		return
	}
	user.ID = uint(uid)

	if err := utils.GetDB().First(&user).Error; err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code: model.DataQuestionCode,
			Msg:  "no data",
		})
		log.Printf("用户:%d不存在,未查找到", user.ID)
		return
	}
	queryUserInfo := model.QueryUserInfo{
		AvatarUrl: user.AvatarUrl,
		ID:        user.ID,
		Nickname:  user.Nickname,
		Signature: user.Signature,
	}
	context.JSON(http.StatusOK, model.Response{
		Code: model.SuccessCode,
		Msg:  "success",
		Data: queryUserInfo,
	})
	log.Printf("成功根据关键字:%s查询到用户", KeywordJSON.Keyword)
}
