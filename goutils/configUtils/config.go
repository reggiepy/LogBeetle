package configUtils

// Config 接口定义通用的配置方法
type Config interface {
	ToJson() (string, error)
	LoadJson(filePath string) error
	ToYaml() (string, error)
	LoadYaml(filePath string) error
}

// JsonConfig 接口定义通用的配置方法
type JsonConfig interface {
	ToJson() (string, error)
	LoadJson(filePath string) error
}

// YamlConfig 接口定义通用的配置方法
type YamlConfig interface {
	ToYaml() (string, error)
	LoadYaml(filePath string) error
}
