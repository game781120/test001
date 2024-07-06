package files

import (
	"bufio"
	"fmt"
	"os"
)

func ReadMarkDown(filePath string) (string, error) {
	// 读取Markdown文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("无法读取文件：%s\n", err)
		return "", err
	}
	return string(content), nil

}

func WriteMarkDown(filePath string, content string) error {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	fmt.Println("Markdown file created successfully.")
	return nil

}

func ReadMarkDownLines(filePath string) ([]string, error) {
	// 打开Markdown文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("无法打开文件：%s\n", err)
		return nil, err
	}
	defer file.Close()

	// 创建Scanner对象
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容
	var lines []string
	for scanner.Scan() {
		fmt.Printf("scanner.Text():%s\n", scanner.Text())
		lines = append(lines, scanner.Text())
	}

	// 检查Scanner是否发生错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件出错：%s\n", err)
		return nil, err
	}

	// 将读取的行拼接成单个字符串
	content := ""
	if len(lines) > 0 {
		content = lines[0]
		for i := 1; i < len(lines); i++ {
			content += "\n" + lines[i]
		}
	}

	return lines, nil
}
