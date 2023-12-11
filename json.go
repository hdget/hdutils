package hdutils

import (
	"encoding/json"
	"github.com/elliotchance/pie/v2"
	"reflect"
)

var (
	EmptyJsonArray  = StringToBytes("[]")
	EmptyJsonObject = StringToBytes("{}")
)

// IsEmptyJsonArray 是否是空json array
func IsEmptyJsonArray(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	return pie.Equals(data, EmptyJsonArray)
}

// IsEmptyJsonObject 是否是空json object
func IsEmptyJsonObject(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	return pie.Equals(data, EmptyJsonObject)
}

// JsonArray 将slice转换成[]byte数据，如果slice为nil或空则返回空json array bytes
func JsonArray(args ...any) []byte {
	if len(args) == 0 || args[0] == nil {
		return EmptyJsonArray
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Slice {
		return EmptyJsonArray
	} else if v.Cap() == 0 {
		return EmptyJsonArray
	}

	jsonData, _ := json.Marshal(args[0])
	return jsonData
}

// JsonObject 将object转换成[]byte数据，如果object为nil或空则返回空json object bytes
func JsonObject(args ...any) []byte {
	if len(args) == 0 || args[0] == nil {
		return EmptyJsonObject
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() == reflect.Pointer {
		v = reflect.ValueOf(v.Elem())
	}

	if v.Kind() != reflect.Struct {
		return EmptyJsonObject
	} else if v.IsZero() {
		return EmptyJsonObject
	}

	jsonData, _ := json.Marshal(args[0])
	return jsonData
}
