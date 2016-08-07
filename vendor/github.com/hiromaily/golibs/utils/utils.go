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
//check type of interface
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

//check type of interface
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

//change string to type
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
// Interface型のString型への変更
func Itos(val interface{}) string {
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

func Itob(val interface{}) bool {
	b, ok := val.(bool)
	if !ok {
		return false
	}
	return b
}

// Interface型のint型への変更
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
	} else {
		return 0
	}
}

// Interface型のuint型への変更
func ItoUi(val interface{}) uint {

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
	} else {
		return 0
	}
}

// Interface型のbyte型->string型への変更
func ItoBS(val interface{}) string {
	if b, ok := val.([]byte); ok {
		return string(b)
	}
	return ""
}

// Interface型のmap[string]int型への変更
func ItoMsi(val interface{}) map[string]int {
	msi, ok := val.(map[string]int)
	if !ok {
		return nil
	}
	return msi
}

// Interface型のtime.Time型への変更
func ItoT(val interface{}) time.Time {
	if t, ok := val.(time.Time); ok {
		return t
	}
	return time.Time{}
}

// Interface型のtime.Time型->string型への変更
func ItoTS(val interface{}) string {
	if t, ok := val.(time.Time); ok {
		return t.String()
	}
	return ""
}

// Interface型をその型を返すfuncを返す
// TODO:型を判別して自動でその型にキャストしたい
// It couldn't be possible.
/*
func ItoAuto(val interface{}) func(interface{}) bool {
	//v := reflect.ValueOf(val).Type()
	//v := reflect.ValueOf(val).Kind()

	v := reflect.ValueOf(val).Kind()

	switch v {
	case reflect.Bool:
		//utils/utils.go:241: cannot use Itob (type func(interface {}) bool) as type func() in return argument
		return Itob
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Itoi
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ItoUi
	case reflect.String:
		return Itos
	default:
		return nil
	}
}
*/

// String型のError型への変更
func Stoe(val string) error {
	return errors.New(val)
}

// String -> Int
func Atoi(str string) (ret int) {
	ret, _ = strconv.Atoi(str)
	return
}

// Int -> String
func Itoa(num int) (ret string) {
	return strconv.Itoa(num)
}

//**********************************************************
// Operate Slice
//**********************************************************
// search string
func SearchString(ary []string, str string) int {

	var retIdx int = -1
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

// TODO:slice
// https://github.com/golang/go/wiki/SliceTricks

//Remove element from end of slice
func Pop(val []interface{}) []interface{} {
	return val[:len(val)-1]
}

//Add element to end of slice
func Push(base []interface{}, val interface{}) []interface{} {
	return append(base, val)
}

//Remove element from first of slice
func Shift(val []interface{}) []interface{} {
	return val[1:]
}

//Add element to first of slice
func Unshift(base []interface{}, val interface{}) []interface{} {
	return append([]interface{}{val}, base...)
}

//change slice data to slice interdace
//https://github.com/golang/go/wiki/InterfaceSlice
func InterfaceSliceInt(dataSlice []int) []interface{} {
	var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
	for i, d := range dataSlice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

func InterfaceSliceMap(dataSlice []map[string]int) []interface{} {
	var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
	for i, d := range dataSlice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

//TODO:work in progress
func GetStructField(strct interface{}) {
	v := reflect.ValueOf(strct)
	fmt.Println(v)
	if v.Kind() != reflect.Struct {
		//return false
		fmt.Println("not struct:", v.Kind())
		return
	}
	for i := 0; i < v.NumField(); i++ {
		//lg.Debugf("v.Field(i).Kind(): %v", v.Field(i).Kind())
		//lg.Debugf("v.Field(i).Pointer(): %v", v.Field(i).Pointer())
		fmt.Printf("v.Field(i).Kind(): %v", v.Field(i).Kind())
		fmt.Printf("v.Field(i).Pointer(): %v", v.Field(i).Pointer())
		//if v.Field(i).Kind() != reflect.Ptr {
		//	return false
		//}
		//if v.Field(0).Pointer() == 0 {
		//	return false
		//}
	}
}

//**********************************************************
// Handle Directory
//**********************************************************
func IsExistDir(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

//**********************************************************
// Handle Error
//**********************************************************
// check error and if so execute panic
func GoPanicWhenError(err error) {
	if err != nil {
		fmt.Println(runtime.Caller(1))
		panic(err)
	}
}

// check error and if so print error
func ShowErrorWhenError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
