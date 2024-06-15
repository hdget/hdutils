package excel

import (
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hdutils/convert"
	"github.com/hdget/hdutils/text"
	"github.com/spf13/cast"
	"strings"
)

type Sheet struct {
	Headers map[int]string
	Rows    []*SheetRow
}

type SheetRow struct {
	Cells []*SheetCell
}

type SheetCell struct {
	RowIndex int
	ColIndex int
	ColName  string
	Value    string
}

func (s *SheetRow) Get(colName string) string {
	index := pie.FindFirstUsing(s.Cells, func(c *SheetCell) bool {
		return strings.EqualFold(c.ColName, text.CleanString(colName))
	})
	if index == -1 {
		return ""
	}
	return text.CleanString(s.Cells[index].Value)
}

func (s *SheetRow) GetInt64(colName string) int64 {
	return cast.ToInt64(s.Get(colName))
}

// GetInt64Slice get comma separated int64 slice
func (s *SheetRow) GetInt64Slice(colName string) []int64 {
	v := s.Get(colName)
	v = strings.ReplaceAll(v, "，", ",")
	return convert.CsvToInt64s(v)
}

// GetStringSlice get comma separated string slice
func (s *SheetRow) GetStringSlice(colName string) []string {
	v := s.Get(colName)
	v = strings.ReplaceAll(v, "，", ",")
	return strings.Split(v, ",")
}
