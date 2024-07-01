package excel

import "github.com/xuri/excelize/v2"

type ExcelWriter interface {
	CreateDefaultSheet(rows []any) error
	CreateSheet(sheetName string, rows []any) error
	Close() error
	GetContent() ([]byte, error)
	Save(filename string) error
}

type excelWriterOption struct {
	defaultSheetName string
	colStyle         *excelize.Style
	cellStyles       map[string]*excelize.Style // key为坐标
}

var (
	defaultExcelWriterOption = excelWriterOption{
		defaultSheetName: "Sheet1",
		colStyle: &excelize.Style{
			NumFmt: 49, // '@'文本占位符。单个@的作用是引用单元格内输入的原始内容，将其以文本格式显示出来,
		},
	}
)

type WriterOption func(*excelWriterOption)

func WithDefaultSheetName(sheetName string) WriterOption {
	return func(option *excelWriterOption) {
		option.defaultSheetName = sheetName
	}
}

func WithColStyle(style *excelize.Style) WriterOption {
	return func(option *excelWriterOption) {
		option.colStyle = style
	}
}

func WithCellStyles(styles map[string]*excelize.Style) WriterOption {
	return func(option *excelWriterOption) {
		option.cellStyles = styles
	}
}
