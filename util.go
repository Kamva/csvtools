package csvtools

import "reflect"

func InterfaceToSlice(l interface{}) []interface{}{
	s := reflect.ValueOf(l)
	if s.Kind() != reflect.Slice {
		panic("InterfaceToSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func ToRecordsSlice(l interface{}) []Record {
	s := reflect.ValueOf(l)
	if s.Kind() != reflect.Slice {
		panic("InterfaceToSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]Record, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface().(Record)
	}

	return ret
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}