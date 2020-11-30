package filesystem

import (
	"context"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"io"
	"log"
	"mime/multipart"
)

type QiniuPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

type QiniuStorage struct {
	config *DriverConfig
}

func NewQiniuStorage(cnf *DriverConfig) *QiniuStorage {
	return  &QiniuStorage{ config: cnf}
}

//七牛云图片上传
func (q *QiniuStorage) Upload(fileInfo *FileInfo, fileHeader *multipart.FileHeader, file multipart.File) (*FileInfo, error) {
	var reader io.Reader = file
	var size = fileHeader.Size

	putPolicy := storage.PutPolicy{
		Scope: q.config.Bucket,
	}
	mac := qbox.NewMac(q.config.AccessKey, q.config.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	formUploader := storage.NewFormUploader(&cfg)

	ret := &QiniuPutRet{} //ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": fileHeader.Filename,
		},
	}

	err := formUploader.Put(context.Background(), &ret, upToken, fileInfo.SaveUrl, reader, size, &putExtra)

	if err != nil {
		log.Println("qiniu put :", err)
		return nil, err
	}

	fileInfo.Hash = ret.Hash
	fileInfo.Size = ret.Fsize

	return fileInfo, nil
}
