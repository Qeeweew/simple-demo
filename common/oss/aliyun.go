package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
	"simple-demo/common/config"
)

func AliyunInit() {
	client, err := oss.New(config.AliyunCfg.Endpoint, config.AliyunCfg.AccessKeyID, config.AliyunCfg.AccessKeySecret)
	if err != nil {
		fmt.Println(err)
		return
	}
	AliyunClient = client
}

func UploadVideoToOss(bucketName string, objectName string, reader multipart.File) (bool, error) {
	bucket, err := AliyunClient.Bucket(bucketName)
	if err != nil {
		return false, err
	}
	err = bucket.PutObject(objectName, reader)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func GetOssVideoUrlAndImgUrl(bucketName string, objectName string) (string, string, error) {
	url := "https://" + bucketName + "." + config.AliyunCfg.Endpoint + "/" + objectName
	return url, url + "?x-oss-process=video/snapshot,t_0,f_jpg,w_0,h_0,m_fast,ar_auto", nil
}
