package filesystem

import (
	"github.com/geiqin/gotools/helper"
	"path"
	"strings"
)

//云储存配置
type CloudConfig struct {
	Driver    string `json:"driver"`     //驱动
	Bucket    string `json:"bucket"`     //Bucket
	AccessKey string `json:"access_key"` //访问key
	SecretKey string `json:"secret_key"` //密钥key
	Transport string `json:"transport"`  //transport
	Domain    string `json:"domain"`     //域名
	Zone      string `json:"zone"`       //区域
	HostUrl   string `json:"host_url"`   //访问url
}

//文件信息
type FileInfo struct {
	Id            int64  `json:"id"`             //对应数据库ID值
	PersistentId  string `json:"persistent_id"`  //持久化ID
	Hash          string `json:"hash"`           //储存Hash值
	Type          string `json:"type"`           //文件类型
	FileName      string `json:"file_name"`      //文件名称
	SavePath      string `json:"save_path"`      //保存路径
	Size          int64  `json:"size"`           //文件大小
	MediaDuration int64  `json:"media_duration"` //媒体文件时长【视频 / 音频】
	MediaWidth    int64  `json:"media_width"`    //媒体文件宽度【图片】
	MediaHeight   int64  `json:"media_height"`   //媒体文件高度【图片】
	Url           string `json:"url"`            //访问URL
}

//获取文件扩展名（如：jpg/png）
func (a FileInfo) GetExtName() string {
	ext := path.Ext(a.FileName)
	extName, _ := helper.Substr(ext, 1, len(ext))
	if extName != "" {
		extName = strings.ToLower(extName)
	}
	return extName
}

//获取文件类型(如：image/video/voice/document/zip/other)
func (a FileInfo) GetType() string {
	ext := a.GetExtName()
	t := "other"
	for k, v := range fileTypes {
		if helper.InArray(v, ext) {
			t = k
			break
		}
	}
	return t
}
