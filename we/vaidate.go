package we

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"reflect"
	"strings"
)

type bindFunc func(obj any) error

func handleErr(err error, data any) error {
	if err == io.EOF {
		return errors.New("请求参数格式错误")
	}
	switch err.(type) {
	case *json.UnmarshalTypeError:
		filedErr := err.(*json.UnmarshalTypeError)
		return errors.New(fmt.Sprintf("%s 字段为%s类型，赋值为%s类型", filedErr.Field, filedErr.Value, filedErr.Type.String()))
	case validator.ValidationErrors:
		errs := err.(validator.ValidationErrors)
		ref := reflect.TypeOf(data)
		for _, fieldError := range errs {
			fieldError.Field()
			snArr := strings.Split(fieldError.StructNamespace(), ".")
			snArr = snArr[1:]
			var filed, _ = ref.Elem().FieldByName(snArr[0])
			snArr = snArr[1:]
			for _, v := range snArr {
				filed, _ = filed.Type.FieldByName(v)
			}
			errTag := fieldError.Tag() + "_msg"
			// 获取对应binding得错误消息
			errTagText := filed.Tag.Get(errTag)
			// 获取统一错误消息
			errText := filed.Tag.Get("msg")
			if errTagText != "" {
				return errors.New(errTagText)
			}
			if errText != "" {
				return errors.New(errText)
			}
			return errors.New(fieldError.Field() + ":" + fieldError.Tag())
		}
	}
	return nil
}
