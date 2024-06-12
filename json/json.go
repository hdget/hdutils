package json

import (
	"encoding/json"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hdutils/convert"
	"reflect"
)

type Utils interface {
	IsEmptyJsonArray(data []byte) bool
	IsEmptyJsonObject(data []byte) bool
	JsonArray(args ...any) []byte
	JsonObject(args ...any) []byte
}

var (
	emptyJsonArray  = convert.StringToBytes("[]")
	emptyJsonObject = convert.StringToBytes("{}")
)

// IsEmptyJsonArray 是否是空json array
func IsEmptyJsonArray(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	return pie.Equals(data, emptyJsonArray)
}

// IsEmptyJsonObject 是否是空json object
func IsEmptyJsonObject(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	return pie.Equals(data, emptyJsonObject)
}

// JsonArray 将slice转换成[]byte数据，如果slice为nil或空则返回空json array bytes
func JsonArray(args ...any) []byte {
	if len(args) == 0 || args[0] == nil {
		return emptyJsonArray
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Slice {
		return emptyJsonArray
	} else if v.Cap() == 0 {
		return emptyJsonArray
	}

	jsonData, _ := json.Marshal(args[0])
	return jsonData
}

// JsonObject 将object转换成[]byte数据，如果object为nil或空则返回空json object bytes
func JsonObject(args ...any) []byte {
	if len(args) == 0 || args[0] == nil {
		return emptyJsonObject
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() == reflect.Pointer {
		v = reflect.ValueOf(v.Elem())
	}

	if v.Kind() != reflect.Struct || v.Kind() != reflect.Map {
		return emptyJsonObject
	} else if v.IsZero() {
		return emptyJsonObject
	}

	jsonData, _ := json.Marshal(args[0])
	return jsonData
}
