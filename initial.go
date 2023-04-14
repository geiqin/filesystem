package filesystem

import "github.com/geiqin/gotools/helper"

var fileTypes map[string][]string

func init() {
	fileTypes = map[string][]string{
		"image":    {"jpg", "jpeg", "png", "gif"},
		"video":    {"avi", "mov", "rmvb", "flv", "mp4", "wmv", "3gp", "rm"},
		"voice":    {"mp3", "wma", "wav"},
		"document": {"txt", "json", "doc", "docx", "xls", "xlsx", ".ppt", ".pptx", "wps", "pdf"},
		"zip":      {"rar", "zip", "gz"},
	}
}

//获得店铺保存路径
func GetStoreSavePath(storeId int64, suffix string) string {
	var flag string
	if storeId > 0 {
		if storeId == 1 {
			flag = "master/"
		} else {
			flag = helper.GetIdentityFlag(storeId, "stores/", "/")
		}
	} else {
		flag = "common/"
	}
	return flag + suffix
}
