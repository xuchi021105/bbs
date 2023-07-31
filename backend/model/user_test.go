package model_test

import (
	"backend/model"
	"backend/utils"
	"testing"
)

func TestGetFansIDs(t *testing.T) {
	utils.InitAndLoadDB()
	user := model.User{}
	user.ID = 2
	var ids []uint
	var err error
	if ids, err = user.GetFansIDs(); err != nil {
		t.Errorf("GetFansIDs函数有问题")
	}
	t.Log("fansIDs:", ids)
}
