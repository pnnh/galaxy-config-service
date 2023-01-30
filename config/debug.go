package config

import (
	"fmt" 
	"os"
)

func LoadDebugConfig(fileName, env string) (string, error) {
	data, err := os.ReadFile("debug/" + fileName) 
	if err != nil {
		return "", fmt.Errorf("读取配置文件出错: %w", err)
	}
	content := string(data)
	return content, nil
}