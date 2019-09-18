/*
  GoLang code created by Jirawat Harnsiriwatanakit https://github.com/kazekim
*/

package pdfinject

import (
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

		value := field
		form[key] = value
	}
	return form
}
