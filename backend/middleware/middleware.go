package middleware

import (
	"backend/model"
	"backend/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 必须登录,有token才能使用的中间件
// 没有token返回一个json,前端来显示要先登录
func JWTAuthMiddleware() func(*gin.Context) {
	return func(context *gin.Context) {
		log.Print("必须要登录鉴权的中间件")
		authHeader := context.Request.Header.Get("Authorization") // 如果没有Authorization的话返回""
		if authHeader == "" {
			context.AbortWithStatusJSON(http.StatusOK, model.UnAuthorizedResponse) // 中间件中必须abort,不然会通过Next()像栈一样并且链式执行下去
			log.Printf("head中没有Authorzation这个字段,未鉴权,请先登录")
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			context.AbortWithStatusJSON(http.StatusOK, model.UnAuthorizedResponse)
			log.Print("没有用Bearer开头,token格式有问题")
			return
		}
		token, err := utils.ParseToken(parts[1])
		if err != nil {
			context.AbortWithStatusJSON(http.StatusOK, model.UnAuthorizedResponse)
			log.Print("token解析失败")
			return
		}
		log.Print("token解析成功,鉴权通过")
		context.Set("myUid", token.Uid) // 将uid信息用Set函数保存到Context上下文中,用Get函数调用
		context.Next()
	}
}

// 不登录,没有token也能使用的中间件
// 没有token可以正常使用,但是有token的话会有额外的效果
func JWTAltAuthMiddleware() func(*gin.Context) {
	return func(context *gin.Context) {
		log.Print("不一定要登录鉴权的中间件")
		authHeader := context.Request.Header.Get("Authorization") // 如果没有Authorization的话返回""
		if authHeader == "" {
			log.Printf("head中没有Authorzation这个字段,走未登录的部分")
			// 这里不用Abort,因为还要继续执行下去的
			// 不在context中用Set函数,如果没有Set就Get的话会得到nil, false(类似断言的返回值),只要对这个返回值做下处理就好了
			context.Next()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			context.AbortWithStatusJSON(http.StatusOK, model.UnAuthorizedResponse)
			log.Print("没有用Bearer开头,token格式有问题")
			return
		}
		token, err := utils.ParseToken(parts[1])
		if err != nil {
			context.AbortWithStatusJSON(http.StatusOK, model.UnAuthorizedResponse)
			log.Print("token解析失败")
			return
		}
		log.Print("token解析成功,鉴权通过")
		context.Set("myUid", token.Uid) // 将uid信息用Set函数保存到Context上下文中,用Get函数调用
		context.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*") // 暂时host
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
