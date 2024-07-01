package excel

type ExcelReader interface {
	ReadSheet(sheetName string) (*Sheet, error)
	ReadAllSheets() ([]*Sheet, error)
}

type excelReaderOption struct {
	headerRowIndex int
}

var (
	defaultExcelReaderOption = excelReaderOption{
		headerRowIndex: 1,
	}
)

type ReaderOption func(*excelReaderOption)

func WithHeaderRowIndex(index int) ReaderOption {
	return func(option *excelReaderOption) {
		option.headerRowIndex = index
	}
}
