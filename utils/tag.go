package utils

import (
	"fmt"
	"reflect"
)

type ParamType int

const (
	form ParamType = iota
	json
)

func (c ParamType) String() string {
	switch c {
	case form:
		return "form"
	case json:
		return "json"
	}
	return ""
}

// BindQuery 绑定Get请求参数
func BindQuery(a int, v interface{}) {
	typeOf := reflect.TypeOf(v)
	valueOf := reflect.ValueOf(v)
	for i := 0; i < typeOf.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := typeOf.Field(i)
		get := fieldType.Tag.Get("form")
		if get != "" {
			t := fieldType.Type
			switch t.Kind() {
			case reflect.Int:
				//s := c.Req.Form.Get(get)
				//intNum, _ := strconv.Atoi(s)
				valueOf.Field(i).Set(reflect.ValueOf(a))
			case reflect.String:

			}
		}
		// 输出成员名和tag
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}
}
