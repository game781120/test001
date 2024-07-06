package files

import (
	"fmt"
	"os"
	"time"
)

func WriteTxtData(fileName, content string) {

	now := time.Now()
	// 格式化日期和时间
	date := now.Format("20060102")
	hour := now.Format("15")
	minute := now.Format("04")
	// 构建文件名
	filename := fmt.Sprintf("%s_%s-%s-%s.txt", fileName, date, hour, minute)
	path := "./temp/"
	// 创建路径
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println("无法创建路径:", err)
		return
	}

	file, err := os.OpenFile(path+filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println(err)
	}
}
