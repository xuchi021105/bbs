package utils

import (
	"strings"
	"testing"
)

func TestGeneToken(t *testing.T) {
	var tokenString string
	var err error
	if tokenString, err = GeneToken(123456); err != nil {
		t.Errorf("生成token的函数: %+v", err)
	}
	t.Logf("tokenString: %s", tokenString)
}

// 解析Bearer token格式
func TestParseToken(t *testing.T) {
	var err error
	var tokenString string
	if tokenString, err = GeneToken(123456); err != nil {
		t.Errorf("生成token的函数: %+v", err)
	}

	parts := strings.SplitN(tokenString, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		t.Error("没有用Bearer开头,token格式有问题")
	}

	var authClaims *AuthClaims
	if authClaims, err = ParseToken(parts[1]); err != nil {
		t.Error("token解析失败")
	}
	t.Logf("authClaims:%+v", *authClaims)
}
