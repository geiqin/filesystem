package filesystem

import (
	"github.com/geiqin/gotools/helper"
	"path"
	"strings"
)

var fileTypes map[string][]string

func init() {
	fileTypes = make(map[string][]string)
	fileTypes["image"] = []string{"jpg", "jpeg", "png", "gif"}
	fileTypes["video"] = []string{"avi", "mov", "rmvb", "flv", "mp4", "wmv", "3gp", "rm"}
	fileTypes["voice"] = []string{"mp3", "wma", "wav"}
	fileTypes["document"] = []string{"txt", "json", "doc", "docx", "xls", "xlsx", ".ppt", ".pptx", "wps", "pdf"}
	fileTypes["zip"] = []string{"rar", "zip", "gz"}
}

func GetExtName(fileName string) string {
	ext := path.Ext(fileName)
	extName, _ := helper.Substr(ext, 1, len(ext))
	if extName != "" {
		extName = strings.ToLower(extName)
	}
	return extName
}

//获取文件类型
func GetFileType(fileName string) string {
	ext := GetExtName(fileName)
	t := ""
	for k, v := range fileTypes {
		if helper.InArray(v, ext) {
			t = k
			break
		}
	}
	return t
}

//验证图片类型
func HasImage(fileName string) bool {
	ext := GetExtName(fileName)
	return helper.InArray(fileTypes["image"], ext)
}

//验证视频类型
func HasVideo(fileName string) bool {
	ext := GetExtName(fileName)
	return helper.InArray(fileTypes["video"], ext)
}

//验证语音类型
func HasVoice(fileName string) bool {
	ext := GetExtName(fileName)
	return helper.InArray(fileTypes["voice"], ext)
}

//验证压缩文件类型
func HasZip(fileName string) bool {
	ext := GetExtName(fileName)
	return helper.InArray(fileTypes["zip"], ext)
}

//验证文档类型
func HasDocument(fileName string) bool {
	ext := GetExtName(fileName)
	return helper.InArray(fileTypes["document"], ext)
}
