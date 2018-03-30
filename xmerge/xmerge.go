package xmerge

import (
	"reflect"

	"github.com/imdario/mergo"
)

func Merge(dst, src interface{}) error {
	if src != nil && !reflect.ValueOf(src).IsNil() {
		return mergo.Merge(dst, src, mergo.WithOverride)
	}

	return nil
}
