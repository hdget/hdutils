package convert

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unsafe"
)

type Utils interface {
	StringToBytes(s string) []byte
	BytesToString(b []byte) string
	ToString(value any) (string, error)
	ToBytes(value any) ([]byte, error)
	CsvToInt64s(strValue string) []int64
	CsvToInt32s(strValue string) []int32
	CsvToInts(strValue string) []int
	Int64sToCsv(int64s []int64) string
	Int32sToCsv(int32s []int32) string
	ToInt64s(strSlice []string) []int64
	ToInt32s(strSlice []string) []int32
	ToInts(strSlice []string) []int
	ToCamelCase(str string) string
	ToSnakeCase(s string) string
	ToSlice(data any) []any
}

var (
	rxCameling = regexp.MustCompile(`[\p{L}\p{N}]+`)
)

// StringToBytes converts text to byte slice without a memory allocation.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to text without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToString 尝试将值转换成字符串
func ToString(value any) (string, error) {
	switch reply := value.(type) {
	case string:
		return reply, nil
	case []byte:
		return BytesToString(reply), nil
	}

	bs, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return BytesToString(bs), nil
}

func ToBytes(value any) ([]byte, error) {
	var data []byte
	switch t := value.(type) {
	case string:
		data = StringToBytes(t)
	case []byte:
		data = t
	default:
		v, err := json.Marshal(value)
		if err != nil {
			return nil, errors.Wrapf(err, "marshal value, value: %v", value)
		}
		data = v
	}
	return data, nil
}

// CsvToInt64s 将逗号分隔的string尝试转换成[1,2,3...]的int64 slice
// Csv means Comma Separated Value
func CsvToInt64s(s string) []int64 {
	if len(s) == 0 {
		return nil
	}

	tokens := strings.Split(s, ",")
	if len(tokens) == 0 {
		return nil
	}

	return ToInt64s(tokens)
}

// CsvToInt32s 将逗号分隔的string尝试转换成[1,2,3...]的int32 slice
// Csv means Comma Separated Value
func CsvToInt32s(s string) []int32 {
	if len(s) == 0 {
		return nil
	}

	tokens := strings.Split(s, ",")
	if len(tokens) == 0 {
		return nil
	}

	return ToInt32s(tokens)
}

// CsvToInts 将逗号分隔的string尝试转换成[1,2,3...]的int slice
// Csv means Comma Separated Value
func CsvToInts(s string) []int {
	if len(s) == 0 {
		return nil
	}

	tokens := strings.Split(s, ",")
	if len(tokens) == 0 {
		return nil
	}

	return ToInts(tokens)
}

// Int64sToCsv 将int64 slice转换成用逗号分隔的字符串: 1,2,3
func Int64sToCsv(int64s []int64) string {
	return strings.Join(cast.ToStringSlice(int64s), ",")
}

// Int32sToCsv 将int32 slice转换成用逗号分隔的字符串: 1,2,3
func Int32sToCsv(int32s []int32) string {
	return strings.Join(cast.ToStringSlice(int32s), ",")
}

// ToInt64s 将string slice转换成[1,2,3...]的int64 slice
func ToInt64s(strSlice []string) []int64 {
	int64s := make([]int64, len(strSlice))
	for i, item := range strSlice {
		int64s[i] = cast.ToInt64(item)
	}
	return int64s
}

// ToInt32s 将string slice转换成[1,2,3...]的int32 slice
func ToInt32s(strSlice []string) []int32 {
	int32s := make([]int32, len(strSlice))
	for i, item := range strSlice {
		int32s[i] = cast.ToInt32(item)
	}
	return int32s
}

// ToInts 将string slice转换成[1,2,3...]的int slice
func ToInts(strSlice []string) []int {
	ints := make([]int, len(strSlice))
	for i, item := range strSlice {
		ints[i] = cast.ToInt(item)
	}
	return ints
}

// ToCamelCase converts from underscore separated form to camel case form.
func ToCamelCase(s string) string {
	byteSrc := []byte(s)
	chunks := rxCameling.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		chunks[idx] = cases.Title(language.English).Bytes(val)
	}
	return string(bytes.Join(chunks, nil))
}

// ToSnakeCase converts from camel case form to underscore separated form.
func ToSnakeCase(s string) string {
	s = ToCamelCase(s)
	runes := []rune(s)
	length := len(runes)
	var out []rune
	for i := 0; i < length; i++ {
		out = append(out, unicode.ToLower(runes[i]))
		if i+1 < length && (unicode.IsUpper(runes[i+1]) && unicode.IsLower(runes[i])) {
			out = append(out, '_')
		}
	}

	return string(out)
}

// ToSlice 将传过来的数据转换成[]any
func ToSlice(data any) []any {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil
	}

	sliceLenth := v.Len()
	sliceData := make([]any, sliceLenth)
	for i := 0; i < sliceLenth; i++ {
		sliceData[i] = v.Index(i).Interface()
	}

	return sliceData
}
