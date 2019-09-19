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
		if st.Type.Kind() == reflect.Bool {
			if field.Bool() {
				value = "Yes"
			}else{
				value = "No"
			}
		}else{
			value = field.String()
		}

		form[key] = value
		fmt.Println(key, " ", value)
	}

	return form
}
