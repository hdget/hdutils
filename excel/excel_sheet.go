package excel

import (
	"github.com/hdget/hdutils/convert"
	"github.com/hdget/hdutils/text"
	"github.com/spf13/cast"
	"strings"
)

type Sheet struct {
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

// GetInt64Slice get comma separated int64 slice
func (r *SheetRow) GetInt64Slice(colName string) []int64 {
	v := r.Get(colName)
	v = strings.ReplaceAll(v, "，", ",")
	return convert.CsvToInt64s(v)
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
