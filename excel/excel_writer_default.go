package excel

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"reflect"
)

type excelWriterImpl struct {
	*excelize.File
	option     *excelWriterOption
	colStyle   int            // 列样式
	cellStyles map[string]int // 单元格样式， axis=>style value
}

const (
	allColAxis = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func NewWriter(options ...WriterOption) (ExcelWriter, error) {
	var option excelWriterOption
	err := copier.Copy(&option, defaultExcelWriterOption)
	if err != nil {
		return nil, err
	}

	writer := &excelWriterImpl{
		File:       excelize.NewFile(),
		option:     &option,
		cellStyles: make(map[string]int),
	}

	for _, apply := range options {
		apply(writer.option)
	}

	writer.colStyle, _ = writer.NewStyle(writer.option.colStyle)
	for axis, style := range writer.option.cellStyles {
		writer.cellStyles[axis], _ = writer.NewStyle(style)
	}

	return writer, nil
}

func (w *excelWriterImpl) CreateDefaultSheet(rows []any) error {
	return w.CreateSheet(w.option.defaultSheetName, rows)
}

func (w *excelWriterImpl) CreateSheet(sheetName string, rows []any) error {
	// new sheet
	sheet, err := w.NewSheet(sheetName)
	if err != nil {
		return errors.Wrap(err, "create sheet")
	}
	w.SetActiveSheet(sheet)

	rowType, err := w.getRowType(rows)
	if err != nil {
		return errors.Wrap(err, "get row type")
	}

	// 通过反射处理数据
	for i := 0; i < rowType.NumField(); i++ {
		// 输出表头
		colName := rowType.Field(i).Tag.Get("col_name")
		if colName == "" {
			continue
		}

		// 如果设置了列坐标，使用列坐标，否则自动生成列坐标
		colAxis := rowType.Field(i).Tag.Get("col_axis")
		if colAxis == "" {
			colAxis = w.genColAxis(i)
		}

		// 设置表头
		headerAxis := fmt.Sprintf("%s%d", colAxis, 1)
		err = w.SetCellValue(sheetName, headerAxis, colName)
		if err != nil {
			return errors.Wrap(err, "generate header")
		}

		// 设置所有有效列的样式
		err = w.SetColStyle(sheetName, colAxis, w.colStyle)
		if err != nil {
			return errors.Wrap(err, "set col colStyle")
		}

		// 输出行数据
		for line, r := range rows {
			lineAxis := fmt.Sprintf("%s%d", colAxis, line+2)
			value := reflect.ValueOf(r).Elem().FieldByName(rowType.Field(i).Name)

			if cellStyle, exist := w.cellStyles[lineAxis]; exist {
				_ = w.SetCellStyle(sheetName, lineAxis, lineAxis, cellStyle)
			}

			err = w.SetCellValue(sheetName, lineAxis, value)
			if err != nil {
				return errors.Wrapf(err, "set cell value, sheetname: %s, axis: %s, value: %v", sheetName, lineAxis, value)
			}
		}

	}

	return nil
}

// Save do nothing
func (w *excelWriterImpl) Save(filename string) error {
	return nil
}

func (w *excelWriterImpl) GetContent() ([]byte, error) {
	// 获取文件写入buffer
	buf, err := w.WriteToBuffer()
	if err != nil {
		return nil, errors.Wrap(err, "excel write to buffer")
	}

	return buf.Bytes(), nil
}

func (w *excelWriterImpl) getRowType(rows []any) (reflect.Type, error) {
	// 通过反射获取表格的标题属性, 生成表格表头
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty rows")
	}

	// 取第一行认为是表头
	firstRowType := reflect.TypeOf(rows[0])
	// 如果是数组或者
	if firstRowType.Kind() == reflect.Ptr || firstRowType.Kind() == reflect.Array || firstRowType.Kind() == reflect.Slice {
		firstRowType = firstRowType.Elem()
	}
	if firstRowType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid struct: %s", firstRowType.String())
	}

	return firstRowType, nil
}

func (w *excelWriterImpl) genColAxis(index int) string {
	mod := index % 26
	s := allColAxis[mod : mod+1]
	var buf bytes.Buffer
	for i := 1; i <= index/26; i++ {
		buf.WriteString("A")
	}
	buf.WriteString(s)
	return buf.String()
}
