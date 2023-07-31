package service

import (
	"fmt"
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	avatarPictureBucket = "bss-avatar-picture"
	endPoint            = "oss-cn-shenzhen.aliyuncs.com"
	OSSServerName       = fmt.Sprintf("http://%s.%s/", avatarPictureBucket, endPoint)
	accessKeyID         = "your accessKey ID"
	accessKeySecret     = "your accessKey Secret"
	client              *oss.Client
	bucket              *oss.Bucket
	err                 error
)

func Init() {
	InitOSSClient()
	InitAvatarPictureBucket()
}

func InitOSSClient() {
	if client, err = oss.New(endPoint, accessKeyID, accessKeySecret); err != nil {
		log.Fatalf("aliyun oss init occur error :%+v", err)
	}
}

func InitAvatarPictureBucket() {

	if bucket, err = client.Bucket(avatarPictureBucket); err != nil {
		log.Fatalf("aliyun oss bss-avatar-picture bucket occur error :%+v", err)
	}

}

func GetOSSClient() *oss.Client {
	return client
}

func GetAvatarPictureBucket() *oss.Bucket {
	return bucket
}
