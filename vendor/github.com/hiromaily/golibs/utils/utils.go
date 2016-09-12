package utils

import (
	"errors"
	"fmt"
	//lg "github.com/hiromaily/golibs/log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

//**********************************************************
// Type Interface
//**********************************************************

// CheckInterface is to check type of interface
func CheckInterface(v interface{}) string {
	//ValueOfでreflect.Value型のオブジェクトを取得
	//v := reflect.ValueOf(val).Type()

	//switch
	switch v.(type) {
	case int, int64, int32, int16, int8:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	case []uint8:
		return "[]uint8"
	//case []byte:
	//	return "[]byte"
	case time.Time:
		return "time.Time"
	default:
		return "default"
	}
}

// CheckInterfaceByIf is to check type of interface
func CheckInterfaceByIf(val interface{}) string {
	//ValueOfでreflect.Value型のオブジェクトを取得
	v := reflect.ValueOf(val).Kind()

	switch v {
	case reflect.Bool:
		return reflect.Bool.String()
	case reflect.Int:
		return reflect.Int.String()
	case reflect.Int8:
		return reflect.Int8.String()
	case reflect.Int16:
		return reflect.Int16.String()
	case reflect.Int32:
		return reflect.Int32.String()
	case reflect.Int64:
		return reflect.Int64.String()
	case reflect.Uint:
		return reflect.Uint.String()
	case reflect.Uint8:
		return reflect.Uint8.String()
	case reflect.Uint16:
		return reflect.Uint16.String()
	case reflect.Uint32:
		return reflect.Uint32.String()
	case reflect.Uint64:
		return reflect.Uint64.String()
	case reflect.Float32:
		return reflect.Float32.String()
	case reflect.Float64:
		return reflect.Float64.String()
	case reflect.Array:
		return reflect.Array.String()
	case reflect.Chan:
		return reflect.Chan.String()
	case reflect.Func:
		return reflect.Func.String()
	case reflect.Interface:
		return reflect.Interface.String()
	case reflect.Map:
		return reflect.Map.String()
	case reflect.Ptr:
		//ptr -> pointer
		return reflect.Ptr.String()
	case reflect.Slice:
		return reflect.Slice.String()
	case reflect.String:
		return reflect.String.String()
	case reflect.Struct:
		return reflect.Struct.String()
	default:
		return ""
	}
}

// StoType is to change string to type
func StoType(typeStr string) reflect.Kind {
	switch typeStr {
	case reflect.Invalid.String():
		return reflect.Invalid
	case reflect.Bool.String():
		return reflect.Bool
	case reflect.Int.String():
		return reflect.Int
	case reflect.Int8.String():
		return reflect.Int8
	case reflect.Int16.String():
		return reflect.Int16
	case reflect.Int32.String():
		return reflect.Int32
	case reflect.Int64.String():
		return reflect.Int64
	case reflect.Uint.String():
		return reflect.Uint
	case reflect.Uint8.String():
		return reflect.Uint8
	case reflect.Uint16.String():
		return reflect.Uint16
	case reflect.Uint32.String():
		return reflect.Uint32
	case reflect.Uint64.String():
		return reflect.Uint64
	case reflect.Uintptr.String():
		return reflect.Uintptr
	case reflect.Float32.String():
		return reflect.Float32
	case reflect.Float64.String():
		return reflect.Float64
	case reflect.Array.String():
		return reflect.Array
	case reflect.Chan.String():
		return reflect.Chan
	case reflect.Func.String():
		return reflect.Func
	case reflect.Interface.String():
		return reflect.Interface
	case reflect.Map.String():
		return reflect.Map
	case reflect.Ptr.String():
		return reflect.Ptr
	case reflect.Slice.String():
		return reflect.Slice
	case reflect.String.String():
		return reflect.String
	case reflect.Struct.String():
		return reflect.Struct
	default:
		return 0
	}
}

//**********************************************************
// Convert type to other type
//**********************************************************

// Itos is to convert interface{} to string
func Itos(val interface{}) string {
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

// Itob is to convert interface{} to bool
func Itob(val interface{}) bool {
	b, ok := val.(bool)
	if !ok {
		return false
	}
	return b
}

// Itoi is to convert interface{} to int
func Itoi(val interface{}) int {

	num64, ok := val.(int64)
	if ok {
		return int(num64)
	}

	num16, ok := val.(int16)
	if ok {
		return int(num16)
	}

	num32, ok := val.(int32)
	if ok {
		return int(num32)
	}

	num, ok := val.(int)
	if ok {
		return int(num)
	}

	return 0
}

// ItoUI is to convert interface{} to uint
func ItoUI(val interface{}) uint {

	num64, ok := val.(uint64)
	if ok {
		return uint(num64)
	}

	num16, ok := val.(uint16)
	if ok {
		return uint(num16)
	}

	num32, ok := val.(uint32)
	if ok {
		return uint(num32)
	}

	num, ok := val.(uint)
	if ok {
		return uint(num)
	}

	return 0
}

// ItoBS is to convert byte[] of interface{} to string
func ItoBS(val interface{}) string {
	if b, ok := val.([]byte); ok {
		return string(b)
	}
	return ""
}

// ItoMsi is to convert map[string] of interface{} to map[string]int
func ItoMsi(val interface{}) map[string]int {
	msi, ok := val.(map[string]int)
	if !ok {
		return nil
	}
	return msi
}

// ItoT is to convert interface{} to time.Time
func ItoT(val interface{}) time.Time {
	if t, ok := val.(time.Time); ok {
		return t
	}
	return time.Time{}
}

// ItoTS is to convert time.Time of interface{} to string
func ItoTS(val interface{}) string {
	if t, ok := val.(time.Time); ok {
		return t.String()
	}
	return ""
}

// Stoe is to convert string to error
func Stoe(val string) error {
	return errors.New(val)
}

// Atoi is to convert string to int
func Atoi(str string) (ret int) {
	ret, _ = strconv.Atoi(str)
	return
}

// Itoa is to convert int to string
func Itoa(num int) (ret string) {
	return strconv.Itoa(num)
}

//**********************************************************
// Operate Slice
// https://github.com/golang/go/wiki/SliceTricks
//**********************************************************

// SearchString is to search string
func SearchString(ary []string, str string) int {

	retIdx := -1
	if len(ary) == 0 {
		return retIdx
	}
	for i, val := range ary {
		if val == str {
			retIdx = i
			break
		}
	}

	return retIdx
}

// Pop is to remove element from end of slice
func Pop(val []interface{}) []interface{} {
	return val[:len(val)-1]
}

// Push is to add element to end of slice
func Push(base []interface{}, val interface{}) []interface{} {
	return append(base, val)
}

// Shift is to remove element from first of slice
func Shift(val []interface{}) []interface{} {
	return val[1:]
}

// Unshift is to add element to first of slice
func Unshift(base []interface{}, val interface{}) []interface{} {
	return append([]interface{}{val}, base...)
}

// SliceIntToInterface is to change slice data of int to slice []interface{}
// https://github.com/golang/go/wiki/InterfaceSlice
func SliceIntToInterface(dataSlice []int) []interface{} {
	interfaceSlice := make([]interface{}, len(dataSlice))
	for i, d := range dataSlice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

// SliceStrToInterface is to change slice data of string to slice []interface{}
func SliceStrToInterface(dataSlice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(dataSlice))
	for i, d := range dataSlice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

// SliceMapToInterface is to change slice data of map[string]int to slice []interface{}
func SliceMapToInterface(dataSlice []map[string]int) []interface{} {
	interfaceSlice := make([]interface{}, len(dataSlice))
	for i, d := range dataSlice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

//**********************************************************
// Handle Directory
//**********************************************************

// IsExistDir is to check existence of directory
func IsExistDir(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

//**********************************************************
// Handle Error
//**********************************************************

// GoPanicWhenError is to execute panic when error
func GoPanicWhenError(err error) {
	if err != nil {
		fmt.Println(runtime.Caller(1))
		panic(err)
	}
}

// ShowErrorWhenError is to show error when error
func ShowErrorWhenError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
