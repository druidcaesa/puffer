package utils

import (
	"fmt"
	"github.com/druidcaesa/puffer"
	"reflect"
	"strconv"
	"strings"
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

// BindQuery Bind Get request parameters
func BindQuery(c *puffer.Context, v interface{}) interface{} {
	return setValueByTag(fmt.Sprintf("%s", form), c, v)
}

// BindJson Body body JSON data submission, data binding method
func BindJson(c *puffer.Context, v interface{}) interface{} {
	return setValueByTag(fmt.Sprintf("%s", json), c, v)
}

//Attribute copy method according to tag
func setValueByTag(tagName string, c *puffer.Context, data interface{}) interface{} {
	// the struct variable
	v := reflect.ValueOf(data).Elem()
	for i := 0; i < v.NumField(); i++ {
		// a reflect.StructField
		fieldInfo := v.Type().Field(i)
		// a reflect.StructTag
		tag := fieldInfo.Tag
		name := tag.Get(tagName)
		//remove possible commas
		name = strings.Split(name, ",")[0]
		if name != "" {
			t := fieldInfo.Type
			switch t.Kind() {
			case reflect.Int:
				get := c.Req.Form.Get(name)
				intNum, _ := strconv.Atoi(get)
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(intNum))
			case reflect.String:
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(c.Req.Form.Get(name)))
			case reflect.Int64:
				parseInt, err := strconv.ParseInt(c.Req.Form.Get(name), 10, 64)
				if err != nil {
					parseInt = 0
				}
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(parseInt))
			}
		}
	}
	return data
}
