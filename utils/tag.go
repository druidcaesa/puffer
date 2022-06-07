package utils

import (
	"fmt"
	"github.com/druidcaesa/puffer"
	"reflect"
	"strconv"
	"strings"
)

type ParamType int

type TagFun interface {
	BindForm(v interface{}) interface{}
	setValueByTag(tagName string, data interface{}) interface{}
	BindJson(v interface{}) interface{}
}

type Tag struct {
	c *puffer.Context
}

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

// BindForm Bind Get request parameters
func (t *Tag) BindForm(v interface{}) interface{} {
	return t.setValueByTag(fmt.Sprintf("%s", json), v)
}

// BindJson Body body JSON data submission, data binding method
func (t *Tag) BindJson(v interface{}) interface{} {
	return t.setValueByTag(fmt.Sprintf("%s", json), v)
}

//Attribute copy method according to tag
func (t *Tag) setValueByTag(tagName string, data interface{}) interface{} {
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
			types := fieldInfo.Type
			switch types.Kind() {
			case reflect.Int:
				get := t.c.Req.Form.Get(name)
				intNum, _ := strconv.Atoi(get)
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(intNum))
			case reflect.String:
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(t.c.Req.Form.Get(name)))
			case reflect.Int64:
				parseInt, err := strconv.ParseInt(t.c.Req.Form.Get(name), 10, 64)
				if err != nil {
					parseInt = 0
				}
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(parseInt))
			}
		}
	}
	return data
}
