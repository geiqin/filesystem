package filesystem

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/geiqin/gotools/helper"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"mime/multipart"
	"strings"
)

//七牛云储存
type QiniuStorage struct {
}

func (q QiniuStorage) GetStorageConfig(conf *CloudConfig) *storage.Config {
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	switch conf.Zone {
	case "huanan":
		cfg.Zone = &storage.ZoneHuanan
	case "huabei":
		cfg.Zone = &storage.ZoneHuabei
	case "huadong":
		cfg.Zone = &storage.ZoneHuadong
	default:
		cfg.Zone = &storage.ZoneHuanan
	}
	return &cfg
}

//上传文件
func (q QiniuStorage) PushFile(ctx context.Context, conf *CloudConfig, fileInfo *FileInfo, fileHeader *multipart.FileHeader, isOverWrite bool) error {
	fileInfo.Size = fileHeader.Size
	fs, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer fs.Close()
	return q.PushData(ctx, conf, fileInfo, fs, isOverWrite)
}

func (q QiniuStorage) PushBytes(ctx context.Context, conf *CloudConfig, fileInfo *FileInfo, data []byte, isOverWrite bool) error {
	buff := bytes.NewBuffer(data)
	fileInfo.Size, _ = helper.ToInt64(buff.Len())
	return q.PushData(ctx, conf, fileInfo, buff, isOverWrite)
}

//获取远程文件是否有效
func (q QiniuStorage) GetRemoteResource(conf *CloudConfig, fileInfo *FileInfo) error {
	q.MakeUrl(conf, fileInfo)
	/*
		_, err := http.Get(fileInfo.HostUrl)
		if err != nil {
			return err
		}

	*/
	return nil
}

//生成完整url
func (q QiniuStorage) MakeUrl(conf *CloudConfig, fileInfo *FileInfo) {
	if !strings.HasSuffix(fileInfo.SavePath, "/") {
		fileInfo.SavePath = fileInfo.SavePath + "/"
	}
	if strings.HasSuffix(conf.HostUrl, "/") {
		fileInfo.Url = conf.HostUrl + fileInfo.SavePath + fileInfo.FileName
	} else {
		fileInfo.Url = conf.HostUrl + "/" + fileInfo.SavePath + fileInfo.FileName
	}
}

//上传数据
func (q QiniuStorage) PushData(ctx context.Context, conf *CloudConfig, fileInfo *FileInfo, reader io.Reader, isOverWrite bool) error {
	if fileInfo == nil {
		return errors.New("fileInfo 不能为空")
	}
	if fileInfo.FileName == "" {
		return errors.New("fileInfo的 FileName 参数不能为空")
	}
	if fileInfo.Size <= 0 {
		return errors.New("fileInfo的 Size 必须大于0")
	}
	q.MakeUrl(conf, fileInfo)
	pathFile := fileInfo.SavePath + fileInfo.FileName
	scope := conf.Bucket
	//是否需要覆盖上传
	if isOverWrite {
		scope = fmt.Sprintf("%s:%s", conf.Bucket, pathFile)
	}
	putPolicy := storage.PutPolicy{
		Scope: scope,
	}
	mac := qbox.NewMac(conf.AccessKey, conf.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	formUploader := storage.NewFormUploader(q.GetStorageConfig(conf))
	ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": fileInfo.FileName,
		},
	}
	err := formUploader.Put(ctx, &ret, upToken, pathFile, reader, fileInfo.Size, &putExtra)

	if err != nil {
		return err
	}
	fileInfo.PersistentId = ret.PersistentID
	fileInfo.Hash = ret.Hash

	return nil
}
