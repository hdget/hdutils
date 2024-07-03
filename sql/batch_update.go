package sql

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
	"text/template"
)

type BatchUpdater interface {
	Add(whenValue, thenValue any) BatchUpdater
	Generate() (string, error)
}

const (
	templateBatchUpdate = `
UPDATE {{.Table}} SET {{.UpdateSet}} = 
CASE
	{{range .Cases}}
	WHEN {{$.WhenColumn}} = {{formatValue .WhenValue}} THEN {{formatValue .ThenValue}}
	{{end}}
END
WHERE {{.WhenColumn}} IN ({{range $index, $element := .Cases}}{{formatValue $element.WhenValue}}{{ if lt $index $.LastIndex }},{{ end }}{{end}});`
)

type mysqlBatchUpdater struct {
	Table      string
	UpdateSet  string
	WhenColumn string
	LastIndex  int
	Cases      []*mysqlBatchUpdateCase
}

type mysqlBatchUpdateCase struct {
	WhenValue any
	ThenValue any
}

func NewMysqlBatchUpdater(table, updateSet, whenColumn string) BatchUpdater {
	return &mysqlBatchUpdater{
		Table:      table,
		UpdateSet:  updateSet,
		WhenColumn: whenColumn,
		Cases:      make([]*mysqlBatchUpdateCase, 0),
	}
}

func (u *mysqlBatchUpdater) Add(whenValue, thenValue any) BatchUpdater {
	u.Cases = append(u.Cases, &mysqlBatchUpdateCase{
		WhenValue: whenValue,
		ThenValue: thenValue,
	})
	return u
}

func (u *mysqlBatchUpdater) Generate() (string, error) {
	if u.Table == "" || u.UpdateSet == "" || u.WhenColumn == "" || len(u.Cases) == 0 {
		return "", errors.New("invalid parameter")
	}
	t, err := template.New("").Funcs(template.FuncMap{
		"formatValue": u.formatValue,
	}).Parse(templateBatchUpdate)
	if err != nil {
		return "", err
	}

	// 获取最后一项的index，方便计算逗号的个数
	u.LastIndex = len(u.Cases) - 1

	// 渲染
	var buf bytes.Buffer
	err = t.Execute(&buf, u)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (u *mysqlBatchUpdater) formatValue(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%.4f", v)
	case []byte:
		return fmt.Sprintf("'%s'", convert.BytesToString(v))
	}
	return fmt.Sprintf("%v", value)
}
