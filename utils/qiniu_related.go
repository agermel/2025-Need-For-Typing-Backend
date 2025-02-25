package utils

import (
	"context"
	"errors"
	"strings"
	"time"
	"type/config"
	"type/models"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

type MusicOrAsset interface {
	models.Song | models.Asset // 限定 T 只能是 Song 或 Asset
}

func qiniuConfig() (*credentials.Credentials, string) {
	accesskey := config.AllConfig.ImageHost.AccessKey
	secretkey := config.AllConfig.ImageHost.SecretKey
	bucket := config.AllConfig.ImageHost.Bucket

	// mac系身份验证工具
	mac := credentials.NewCredentials(accesskey, secretkey)
	return mac, bucket
}

func GetToken() (string, error) {
	mac, bucket := qiniuConfig()
	// 生成创建策略，包含有效时间和存储对象
	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}

	// 获取上传凭证
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return upToken, err
	}
	return upToken, nil
}

func GetList[T MusicOrAsset](key string) ([]T, error) {
	mac, bucketName := qiniuConfig()
	var results []T

	objectsManager := objects.NewObjectsManager(
		&objects.ObjectsManagerOptions{
			Options: http_client.Options{Credentials: mac},
		},
	)
	bucket := objectsManager.Bucket(bucketName)

	keyPrefix := key + "/"

	iter := bucket.List(
		context.Background(),
		&objects.ListObjectsOptions{
			Prefix: keyPrefix,
		},
	)
	defer iter.Close()

	var objectInfo objects.ObjectDetails
	for iter.Next(&objectInfo) {
		var item any
		switch key {
		case "music":
			item = models.Song{
				File_id: config.AllConfig.ImageHost.Domain + "/" + objectInfo.Name,
				Title:   strings.TrimPrefix(objectInfo.Name, keyPrefix),
			}
		case "asset":
			item = models.Asset{
				File_id: config.AllConfig.ImageHost.Domain + "/" + objectInfo.Name,
				Name:    strings.TrimPrefix(objectInfo.Name, keyPrefix),
			}
		default:
			return nil, errors.New("不支持的 key 类型")
		}

		results = append(results, item.(T)) // 这里的转换更安全
	}

	return results, nil
}
