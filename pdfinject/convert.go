/*
  GoLang code created by Jirawat Harnsiriwatanakit https://github.com/kazekim
*/

package pdfinject

import (
	"fmt"
	"reflect"
	"strconv"
)

func structToForm(data interface{}) map[string]interface{} {

	form := make(map[string]interface{})
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	for i := 0; i < v.NumField(); i++ {

		field := v.Field(i)
		st := t.Field(i)
		key := st.Name

		tag := st.Tag.Get("pdf")
		if tag != "" {
			key = tag
		}

		var value string

		switch st.Type.Kind() {
		case reflect.Bool:
			if field.Bool() {
				value = "Yes"
			} else {
				value = "No"
			}
			form[key] = value

		case reflect.Slice, reflect.Array:
			for j:= 0; j < field.Len(); j++ {

				sf := field.Index(j)
				sst := sf.Type()
				for k := 0; k < sf.NumField(); k++ {
					key = fmt.Sprint(sst.Field(k).Name,j+1)
					value = sf.Field(k).String()
					form[key] = value

				}
			}
		default:
			value = field.String()
			form[key] = value
			
		}
	}

	return form
}

func convertMapValue(k string, v interface{}, out *map[string]interface{}) {
	o := *out
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Bool:
		if reflect.ValueOf(v).Bool() {
			o[k] = "Yes"
		} else {
			o[k] = "No"
		}
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(v)
		for i:= 0; i < s.Len(); i++ {

			si := s.Index(i)
			sit := reflect.TypeOf(si.Interface())
			siv := reflect.ValueOf(si.Interface())

			if sit.Kind() == reflect.Map {
				for _,key := range siv.MapKeys() {
					convertMapValue(key.String()+strconv.Itoa(i+1),siv.MapIndex(key).Interface(),&o)
				}
			}
		}
	default:
		o[k] = reflect.ValueOf(v).String()
	}
}

func prepareMap(m map[string]interface{}) *map[string]interface{}{

	result := make(map[string]interface{})

	for k,v := range m {

		convertMapValue(k,v,&result)
	}

	return &result
}
