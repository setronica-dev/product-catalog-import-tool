package utils

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

var floatType = reflect.TypeOf(float64(0))
var stringType = reflect.TypeOf("")

func GetFloat(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		v := reflect.ValueOf(unk)
		v = reflect.Indirect(v)
		if v.Type().ConvertibleTo(floatType) {
			fv := v.Convert(floatType)
			return fv.Float(), nil
		} else if v.Type().ConvertibleTo(stringType) {
			sv := v.Convert(stringType)
			s := sv.String()
			return strconv.ParseFloat(s, 64)
		} else {
			return math.NaN(), fmt.Errorf("can't convert %v to float64", v.Type())
		}
	}
}

func InArray(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return exists, index
		}
	}

	return exists, index
}

func ArrayToString(A []string, delimiter string) string {

	var buffer bytes.Buffer
	for i := 0; i < len(A); i++ {
		buffer.WriteString(A[i])
		if i != len(A)-1 {
			buffer.WriteString(delimiter)
		}
	}

	return buffer.String()
}

func MapInterfaceToSliceStrings(item map[string]interface{}) []string {
	slice := make([]string, 0)
	for _, v := range item {
		slice = append(slice, fmt.Sprintf("%v", v))
	}
	return slice
}

func MapInterfaceKeysToSliceStrings(item map[string]interface{}) []string {
	slice := make([]string, 0)
	for k, _ := range item {
		slice = append(slice, k)
	}
	return slice
}

func SliceDiff(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func RevertMapKeyValue(source map[string]string) map[string]string {
	index := make(map[string]string, len(source))
	for i, v := range source {
		index[v] = i
	}
	return index
}

func GetMapOrDefault(key string, source map[string]string) string {
	if val, ok := source[key]; ok {
		return val
	}
	return key
}
