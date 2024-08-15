package jsonUtils

import (
	"encoding/json"
	"fmt"
)

func AnyToJson(data interface{}, format string) (string, error) {
	var jsonData []byte
	var err error
	switch format {
	case "simple":
		// 将结构体序列化为 JSON
		jsonData, err = json.Marshal(data)
	case "humanReadable":
		// 将结构体序列化为 JSON
		jsonData, err = json.MarshalIndent(data, "", "    ")
	default:
		return "", fmt.Errorf("不支持的format类型: %s", format)
	}
	if err != nil {
		return "", fmt.Errorf("配置信息json序列化失败: %v", err)
	}
	return string(jsonData), nil
}
