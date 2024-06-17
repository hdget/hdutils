package excel

import (
	"github.com/hdget/hdutils/text"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"net/http"
)

type httpExcelReader struct {
	*excelize.File
	option *excelReaderOption
	Sheets map[string]*Sheet
}

func NewHttpReader(url string, options ...ReaderOption) (ExcelReader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 读取需要处理的源excel文件
	f, err := excelize.OpenReader(resp.Body)
	if err != nil {
		return nil, err
	}

	reader := &httpExcelReader{File: f, option: defaultExcelReaderOption}
	for _, apply := range options {
		apply(reader.option)
	}
	return reader, nil
}

func (r httpExcelReader) ReadSheet(sheetName string) (*Sheet, error) {
	rows, err := r.Rows(sheetName)
	if err != nil {
		return nil, errors.Wrapf(err, "read rows, sheet: %s", sheetName)
	}

	// 读取表头, 跳过r.option.headerRowIndex行
	for i := 0; i < r.option.headerRowIndex; i++ {
		rows.Next()
	}

	headerRow, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrapf(err, "read header, sheet: %s", sheetName)
	}

	headerIndexes := make(map[string]int)
	headers := make([]string, len(headerRow))
	for i, colCell := range headerRow {
		header := text.CleanString(colCell)
		headerIndexes[header] = i
		headers[i] = header
	}

	s := &Sheet{HeaderIndexes: headerIndexes, Headers: headers, Rows: make([]*SheetRow, 0)}

	// 读取数据
	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return nil, errors.Wrapf(err, "read data, sheet: %s", sheetName)
		}

		line := &SheetRow{Sheet: s, Columns: make([]string, 0)}
		line.Columns = append(line.Columns, columns...)
		
		s.Rows = append(s.Rows, line)
	}

	return s, nil
}

func (r httpExcelReader) ReadAllSheets() ([]*Sheet, error) {
	sheets := make([]*Sheet, 0)
	for _, sheetName := range r.GetSheetList() {
		sheet, err := r.ReadSheet(sheetName)
		if err != nil {
			return nil, errors.Wrapf(err, "read sheet, sheet: %s", sheetName)
		}

		sheets = append(sheets, sheet)
	}
	return sheets, nil
}
