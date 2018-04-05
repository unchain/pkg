package xdeep

import (
	"reflect"

	"github.com/imdario/mergo"
)

func New(obj interface{}) interface{} {
	return deepNew(obj).Interface()
}

func deepNew(obj interface{}) reflect.Value {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Ptr:
		return deepNewPtr(obj)
	default:
		return reflect.New(reflect.TypeOf(obj)).Elem()
	}
}

func deepNewPtr(obj interface{}) reflect.Value {
	res := reflect.New(reflect.TypeOf(obj).Elem()).Elem()

	for i := 0; i < res.NumField(); i++ {
		field := res.Field(i)
		//spew.Dump("field: ", field)
		newField := deepNew(field.Interface())
		//spew.Dump(newField)

		field.Set(newField)
	}

	return res.Addr()
}

func Merge(dst, src interface{}) error {
	if src != nil && !reflect.ValueOf(src).IsNil() {
		return mergo.Merge(dst, src, mergo.WithOverride)
	}

	return nil
}
