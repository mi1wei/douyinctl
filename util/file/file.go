package file

import (
	"fmt"
	"os"
)

func IsExists(platform string, userId string) bool {
	cookieFilePath := fmt.Sprintf("./.%s/cookies_%s.json", platform, userId)
	_, err := os.Stat(cookieFilePath)

	// 判断文件是否存在
	if err != nil {
		return false
	}
	return true
}

func Read(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	for {
		var line string
		_, err := fmt.Fscanln(file, &line)
		if err != nil {
			break
		}
		lines = append(lines, line)
	}

	return lines, nil
}
