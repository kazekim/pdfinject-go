/*
  GoLang code created by Jirawat Harnsiriwatanakit https://github.com/kazekim
*/

package pdfinject

import (
	"fmt"
	"reflect"
)

func structToForm(data interface{}) Form {

	form := make(Form)
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
					key = fmt.Sprint(sst.Field(k).Name,j)
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
