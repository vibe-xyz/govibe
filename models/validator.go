package models

import (
	"fmt"

	"github.com/go-playground/validator"
)

var (
	validate = validator.New()
)

func ValidateStruct(val interface{}) (err error) {
	err = validate.Struct(val)
	if nil != err {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return
		}
		var str string
		for _, err2 := range err.(validator.ValidationErrors) {
			tmp := fmt.Sprintf("%s %s %s [%v]",
				err2.StructField(), err2.Tag(), err2.Param(), err2.Value())
			if len(str) > 0 {
				str += " | "
			}
			str += tmp
		}
		err = fmt.Errorf(str)
	}
	return
}
