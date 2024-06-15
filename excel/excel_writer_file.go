package excel

import (
	"github.com/pkg/errors"
)

type fileExcelWriterImpl struct {
	*excelWriterImpl
}

func NewFileWriter(options ...WriterOption) (ExcelWriter, error) {
	w, err := NewWriter(options...)
	if err != nil {
		return nil, err
	}

	return &fileExcelWriterImpl{
		excelWriterImpl: w.(*excelWriterImpl),
	}, nil
}

func (w *fileExcelWriterImpl) Save(filename string) error {
	err := w.SaveAs(filename)
	if err != nil {
		return errors.Wrap(err, "save excel file")
	}
	return nil
}
