package libs

import "strings"

// 页面编写Case时，传入的json参数有可能带有换行、tab格式
func HandleJson(json string) (compress string) {
	// 去除首尾空格
	json = strings.TrimSpace(json)
	// 去除\r
	json = strings.ReplaceAll(json, "\r", "")
	// 去除\n
	json = strings.ReplaceAll(json, "\n", "")
	// 去除\t
	compress = strings.ReplaceAll(json, "\t", "")
	// 去除'冒号'左侧空格
	leftflag := strings.Contains(compress, " :")
	for leftflag {
		compress = strings.ReplaceAll(compress, " :", ":")
		leftflag = strings.Contains(compress, " :")
	}
	// 去除'冒号'右侧空格
	rightflag := strings.Contains(compress, ": ")
	for rightflag {
		compress = strings.ReplaceAll(compress, ": ", ":")
		rightflag = strings.Contains(compress, ": ")
	}

	return
}
