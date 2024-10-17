package excel

import (
	"github.com/hdget/hdutils/convert"
	"github.com/hdget/hdutils/text"
	"github.com/spf13/cast"
	"strings"
)

type Sheet struct {
	Name          string
	HeaderIndexes map[string]int // headerName => index
	Headers       []string       // headers
	Rows          []*SheetRow
}

type SheetRow struct {
	Sheet   *Sheet
	Columns []string
}

func (r *SheetRow) Get(colName string) string {
	index, exists := r.Sheet.HeaderIndexes[colName]
	if !exists {
		return ""
	}

	// 检查是否越界
	if index > len(r.Columns)-1 {
		return ""
	}

	return text.CleanString(r.Columns[index])
}

func (r *SheetRow) GetInt64(colName string) int64 {
	return cast.ToInt64(r.Get(colName))
}

func (r *SheetRow) GetInt32(colName string) int32 {
	return cast.ToInt32(r.Get(colName))
}

func (r *SheetRow) GetInt(colName string) int {
	return cast.ToInt(r.Get(colName))
}

func (r *SheetRow) GetFloat64(colName string) float64 {
	return cast.ToFloat64(r.Get(colName))
}

// GetInt64Slice get comma separated int64 slice
func (r *SheetRow) GetInt64Slice(colName string) []int64 {
	v := r.Get(colName)
	v = strings.ReplaceAll(v, "，", ",")
	return convert.CsvToInt64s(v)
}

// GetInt32Slice get comma separated int32 slice
func (r *SheetRow) GetInt32Slice(colName string) []int32 {
	v := r.Get(colName)
	v = strings.ReplaceAll(v, "，", ",")
	return convert.CsvToInt32s(v)
}

// GetIntSlice get comma separated int slice
func (r *SheetRow) GetIntSlice(colName string) []int {
	v := r.Get(colName)
	v = strings.ReplaceAll(v, "，", ",")
	return convert.CsvToInts(v)
}

// GetStringSlice get comma separated string slice
func (r *SheetRow) GetStringSlice(colName string) []string {
	v := r.Get(colName)
	if v == "" {
		return nil
	}

	v = strings.ReplaceAll(v, "，", ",")
	return strings.Split(v, ",")
}
