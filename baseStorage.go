package filesystem

import (
	"mime/multipart"
	"strings"
)

//基类
type BaseStorage struct {
	mediaHandle MediaHandle
}

//获取远程文件是否有效
func (b BaseStorage) GetRemoteResource(conf *CloudConfig, fileInfo *FileInfo) error {
	b.MakeUrl(conf, fileInfo)
	/*
		_, err := http.Get(fileInfo.HostUrl)
		if err != nil {
			return err
		}
	*/
	return nil
}

//生成完整url
func (b BaseStorage) MakeUrl(conf *CloudConfig, fileInfo *FileInfo) {
	if !strings.HasSuffix(fileInfo.SavePath, "/") {
		fileInfo.SavePath = fileInfo.SavePath + "/"
	}
	if strings.HasSuffix(conf.HostUrl, "/") {
		fileInfo.Url = conf.HostUrl + fileInfo.SavePath + fileInfo.FileName
	} else {
		fileInfo.Url = conf.HostUrl + "/" + fileInfo.SavePath + fileInfo.FileName
	}
}

//处理多媒体文件数据【图片/视频/音频】
func (b BaseStorage) MakeMediaData(fileInfo *FileInfo, file *multipart.FileHeader) error {
	var err error
	fileType := fileInfo.GetType()
	if fileType == "image" {
		fileInfo.MediaWidth, fileInfo.MediaHeight, err = b.mediaHandle.GetImageScale(file)
	} else if fileType == "video" || fileType == "voice" {
		fileInfo.MediaDuration, err = b.mediaHandle.GetVideoAndVoiceDuration(file)
	}
	return err
}
