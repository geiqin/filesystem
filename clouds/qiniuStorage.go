package clouds

import (
	"context"
	"github.com/geiqin/filesystem"
	"github.com/geiqin/xconfig/model"
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
	driverConf *model.FileSystemInfo
}

func NewQiniuStorage(cnf *model.FileSystemInfo) *QiniuStorage {
	return  &QiniuStorage{ driverConf: cnf}
}

//七牛云图片上传
func (q *QiniuStorage) Upload(fileInfo *filesystem.FileInfo, fileHeader *multipart.FileHeader, file multipart.File) (*filesystem.FileInfo, error) {
	var reader io.Reader = file
	var size = fileHeader.Size

	putPolicy := storage.PutPolicy{
		Scope: q.driverConf.Bucket,
	}
	mac := qbox.NewMac(q.driverConf.AccessKey, q.driverConf.SecretKey)
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