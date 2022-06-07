package utils

import (
	json2 "encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type ParamType int

type TagFun interface {
	BindForm(v interface{}) (bool, error)
	setValueByTag(tagName string, data interface{}) (bool, error)
	BindJson(v interface{}) (bool, error)
	setValueBody(data interface{}) (bool, error)
}

type Tag struct {
	R *http.Request
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
func (t *Tag) BindForm(v interface{}) (bool, error) {
	return t.setValueByTag(fmt.Sprintf("%s", form), v)
}

// BindJson Body body JSON data submission, data binding method
func (t *Tag) BindJson(v interface{}) (bool, error) {
	return t.setValueByTag(fmt.Sprintf("%s", json), v)
}

//Attribute copy method according to tag
func (t *Tag) setValueByTag(tagName string, data interface{}) (bool, error) {
	if tagName == fmt.Sprintf("%s", json) {
		return t.setValueBody(data)
	}
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
				get := t.R.URL.Query().Get(name)
				intNum, err := strconv.Atoi(get)
				if err != nil {
					return false, err
				}
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(intNum))
			case reflect.String:
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(t.R.URL.Query().Get(name)))
			case reflect.Int64:
				parseInt, err := strconv.ParseInt(t.R.URL.Query().Get(name), 10, 64)
				if err != nil {
					return false, err
				}
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(parseInt))
			}
		}
	}
	return true, nil
}

func (t *Tag) setValueBody(data interface{}) (bool, error) {
	err := json2.NewDecoder(t.R.Body).Decode(data)
	if err != nil {
		return false, err
	}
	return true, nil
}
