package configUtils

import "gopkg.in/yaml.v3"

type YamlConfigBase struct{}

// ToYaml 将配置转换为YAML格式
func (yc *YamlConfigBase) ToYaml(c interface{}) string {
	yamlData, _ := yaml.Marshal(c)
	return string(yamlData)
}

// LoadYaml 从YAML文件加载配置
func (yc *YamlConfigBase) LoadYaml(c interface{}, data string) error {
	if err := yaml.Unmarshal([]byte(data), c); err != nil {
		return err
	}
	return nil
}
