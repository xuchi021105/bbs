package http

import (
	"backend/handler"
	"backend/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func router(engine *gin.Engine) {

	log.Print("启动gin服务")

	// 因为axios的原因,导致有些GET方法变成POST方法(还有DELTE方法,变成了PUT),但是在apifox的文档里面方法不变
	engine.POST("/login", handler.Login)
	engine.POST("/register", handler.Register)

	userGroup := engine.Group("/user")
	userGroup.POST("/followers", handler.GetFollowers)
	userGroup.POST("/followMePersonList", handler.GetFollowMePersonList)
	userAuthGroup := userGroup.Group("/auth", middleware.JWTAuthMiddleware())
	userAuthGroup.POST("/followUid", handler.FollowUser)
	userAuthGroup.PUT("/followUid", handler.CancelFollowedUser)
	userAuthGroup.PUT("/info", handler.UpdateUserInfo)
	userAuthGroup.POST("/avatarPicture", handler.UploadAvatarPicture)
	userAltAuthGroup := userGroup.Group("/altAuth", middleware.JWTAltAuthMiddleware())
	userAltAuthGroup.POST("/info", handler.GetUserInfoByUid)

	articleGroup := engine.Group("/article")
	articleAuthGroup := articleGroup.Group("/auth", middleware.JWTAuthMiddleware())
	articleAuthGroup.POST("/article", handler.PostArticle)
	articleAuthGroup.PUT("/article", handler.UpdateArticle)
	articleAuthGroup.DELETE("/article", handler.DeleteArticle)
	articleAuthGroup.POST("/star", handler.StarArticle)
	articleAuthGroup.PUT("/star", handler.CancelStaredArticle)
	articleAltAuthGroup := articleGroup.Group("/altAuth", middleware.JWTAltAuthMiddleware())
	articleAltAuthGroup.POST("/getArticle", handler.GetArticleInfo)
	articleAltAuthGroup.POST("/articles", handler.GetUserAllArticles)
	articleAltAuthGroup.POST("/userStaredArticles", handler.GetUserAllStaredArticles)

	commentGroup := engine.Group("/comment")
	commentGroup.POST("/getCommentArticle", handler.GetAllCommentsAndReplysInArticle)
	commentAuthGroup := commentGroup.Group("/auth", middleware.JWTAuthMiddleware())
	commentAuthGroup.POST("/commentArticle", handler.CommentArticle)

	replyGroup := engine.Group("/reply")
	replyAuthGroup := replyGroup.Group("/auth", middleware.JWTAuthMiddleware())
	replyAuthGroup.POST("/replyComment", handler.ReplyComment)

	queryGroup := engine.Group("/query")
	queryGroup.POST("/articleByKeyword", handler.QueryArticleByKeyword)
	queryGroup.POST("/userByKeyword", handler.QueryUserByKeyword)
	queryGroup.POST("/userByUid", handler.QueryUserByUID)
}

func Run(port string) {

	engine := gin.Default()

	// 不用这个中间件,对OPTIONS方法实现不了,换一个方法实现CORS
	// engine.Use(cors.Default()) // 使用cors跨域中间件,并且允许所有域名访问(*)

	engine.Use(middleware.Cors())

	router(engine)
	if err := engine.Run(port); err != nil {
		log.Fatalf("gin框架Run时启动出错: %+v", err)
	}

	log.Printf("service started on port %s", port)
}
