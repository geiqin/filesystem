package filesystem

import (
	"context"
	"github.com/geiqin/xconfig/model"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"log"
	"mime/multipart"
)

type QiniuStorage struct {
	driverConf *model.FileSystemInfo
	bucket     string
}

func NewQiniuStorage(cnf *model.FileSystemInfo) *QiniuStorage {
	return &QiniuStorage{driverConf: cnf, bucket: cnf.Bucket}
}

//分片上传
func (q *QiniuStorage) UploadBySlices(fileInfo *FileInfo, fileHeader *multipart.FileHeader, file multipart.File) (*FileInfo, error) {
	var reader io.ReaderAt = file
	var size = fileHeader.Size

	putPolicy := storage.PutPolicy{
		Scope: q.driverConf.Bucket,
	}
	mac := qbox.NewMac(q.driverConf.AccessKey, q.driverConf.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	ret := storage.PutRet{}

	putExtra := storage.RputV2Extra{}
	uploader := storage.NewResumeUploaderV2(q.getCfg())
	err := uploader.Put(context.Background(), &ret, upToken, fileInfo.SaveUrl, reader, size, &putExtra)
	if err != nil {
		log.Println("qiniu put :", err)
		return nil, err
	}
	ss := uploader.InitParts()
	ss.
		fileInfo.PersistentId = ret.PersistentID
	fileInfo.Hash = ret.Hash
	fileInfo.Size = size

	return fileInfo, nil
}

//单文件上传
func (q *QiniuStorage) Upload(fileInfo *FileInfo, fileHeader *multipart.FileHeader, file multipart.File) (*FileInfo, error) {
	var reader io.Reader = file
	var size = fileHeader.Size

	putPolicy := storage.PutPolicy{
		Scope: q.driverConf.Bucket,
	}
	mac := qbox.NewMac(q.driverConf.AccessKey, q.driverConf.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	formUploader := storage.NewFormUploader(q.getCfg())
	ret := storage.PutRet{}

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
	fileInfo.PersistentId = ret.PersistentID
	fileInfo.Hash = ret.Hash
	fileInfo.Size = size

	return fileInfo, nil
}

func (q QiniuStorage) getCfg() *storage.Config {
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	return &cfg
}

func (q QiniuStorage) bucketManager() *storage.BucketManager {
	mac := qbox.NewMac(q.driverConf.AccessKey, q.driverConf.SecretKey)
	bucketManager := storage.NewBucketManager(mac, q.getCfg())
	return bucketManager
}

//删除文件
func (q QiniuStorage) Delete(fileKey string, bucket ...string) bool {
	bucketManager := q.bucketManager()
	err := bucketManager.Delete(q.bucket, fileKey)
	if err != nil {
		log.Println("qiniu delete err:", err)
		return false
	}
	return true
}
