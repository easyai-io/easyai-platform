package structure

import (
	"github.com/jinzhu/copier"
)

// Copy 结构体映射
func Copy(fromValue, toValue interface{}) error {
	return copier.Copy(toValue, fromValue)
}
