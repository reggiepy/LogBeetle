package arrayUtils

import (
	"fmt"
	"reflect"
)

// InArray 元素是否在数组(切片/字典)内.
func InArray(needle interface{}, arr interface{}) (r bool) {
	val := reflect.ValueOf(arr)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				r = true
				return
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				r = true
				return
			}
		}
	default:
		panic(fmt.Errorf("[InArray]arr type must be array, slice or map").(any))
	}
	return
}
