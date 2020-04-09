package misc

import "reflect"

// InterfaceSlice will convert any slice to interface slice
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// Field2Tag for get field to tag map
func Field2Tag(s interface{}, tagName string) (m map[string]string) {
	m = make(map[string]string)

	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		m[v.Field(i).Type().Name()] = tag
	}

	return
}

const (
	// MinUint for min uint
	MinUint uint = 0
	// MaxUint for max unit
	MaxUint = ^MinUint
	// MaxInt for max int
	MaxInt = int(MaxUint >> 1)
	// MinInt for min int
	MinInt = ^MaxInt
)
