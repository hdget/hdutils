package excel

import (
	"fmt"
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

func NewWriter(options ...WriterOption) (ExcelWriter, error) {
	writer := &excelWriterImpl{
		File:       excelize.NewFile(),
		option:     defaultExcelWriterOption,
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
		colAxis := rowType.Field(i).Tag.Get("col_axis")

		if colName != "" {
			axis := fmt.Sprintf("%s%d", colAxis, 1)

			err = w.SetCellValue(sheetName, axis, colName)
			if err != nil {
				return errors.Wrap(err, "generate header")
			}
		}

		// 设置所有有效列的样式
		if colAxis != "" {
			err = w.SetColStyle(sheetName, colAxis, w.colStyle)
			if err != nil {
				return errors.Wrap(err, "set col colStyle")
			}
		}

		// 如果有列坐标的输出行数据
		if colAxis != "" {
			for line, r := range rows {
				axis := fmt.Sprintf("%s%d", colAxis, line+2)
				value := reflect.ValueOf(r).Elem().FieldByName(rowType.Field(i).Name)

				if cellStyle, exist := w.cellStyles[axis]; exist {
					_ = w.SetCellStyle(sheetName, axis, axis, cellStyle)
				}

				err = w.SetCellValue(sheetName, axis, value)
				if err != nil {
					return errors.Wrapf(err, "set cell value, sheetname: %s, axis: %s, value: %v", sheetName, axis, value)
				}
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
