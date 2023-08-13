package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

var AliyunClient *oss.Client

func Init() {
	AliyunInit()
}
