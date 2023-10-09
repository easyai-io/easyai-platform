package yaml

import (
	"gopkg.in/yaml.v3"
)

// 定义YAML操作
var (
	Marshal    = yaml.Marshal
	Unmarshal  = yaml.Unmarshal
	NewDecoder = yaml.NewDecoder
	NewEncoder = yaml.NewEncoder
)

// MarshalToString YAML编码为字符串
func MarshalToString(v interface{}) string {
	s, err := Marshal(v)
	if err != nil {
		return ""
	}
	return string(s)
}
