package main

import (
	"backend/config"
	"backend/http"
	"backend/model"
	"backend/service"
	"backend/utils"
)

func main() {

	config.InitConfig()
	utils.InitAndLoadDB()
	service.Init()
	db := utils.GetDB()
	db.AutoMigrate(&model.User{}, &model.Article{}, &model.Comment{}, &model.Reply{})
	http.Run(":" + model.ServerPort)
}
