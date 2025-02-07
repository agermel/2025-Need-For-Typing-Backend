package service

import (
	"context"
	"io"

	"type/config"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

func GetImageURL(key string) string {
	// 拼接素材URL
	return config.AllConfig.ImageHost.Domain + "/" + key
}

// 上传文件到七牛云 参数：本地路径，在存储时的名称路径，数据流
func UploadToQiniu(key string, reader io.Reader) (string, error) {
	// 构建身份验证工具
	mac := credentials.NewCredentials(config.AllConfig.ImageHost.AccessKey, config.AllConfig.ImageHost.SecretKey)

	// 上传管理器
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials:         mac,
			AccelerateUploading: true,
		},
	})
	err := uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: config.AllConfig.ImageHost.Bucket, // 大文件夹
		ObjectName: &key,                              // 上传文件文件名
		FileName:   "Togawa Cooperation",
		// CustomVar待设定
	}, nil)
	if err != nil {
		return "", err
	}

	return GetImageURL(key), err
}
