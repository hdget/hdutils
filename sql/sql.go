package sql

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
	"text/template"
)

type BatchUpdateCaseWhen struct {
	WhenValue any
	ThenValue any
}

const (
	templateBatchUpdate = `
UPDATE {{.Table}} SET {{.UpdateColumn}} = 
CASE
	{{range .Cases}}
	WHEN {{$.WhenColumn}} = {{formatValue .WhenValue}} THEN {{formatValue .ThenValue}}
	{{end}}
END
WHERE {{.WhenColumn}} IN ({{range $index, $element := .Cases}}{{formatValue $element.WhenValue}}{{ if lt $index $.LastIndex }},{{ end }}{{end}});`
)

func BatchUpdate(table, updateColumn, whenColumn string, caseWhens []*BatchUpdateCaseWhen) (string, error) {
	if table == "" || updateColumn == "" || whenColumn == "" || len(caseWhens) == 0 {
		return "", errors.New("invalid parameter")
	}
	t, err := template.New("").Funcs(template.FuncMap{
		"formatValue": formatValue,
	}).Parse(templateBatchUpdate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, map[string]interface{}{
		"Table":        table,
		"UpdateColumn": updateColumn,
		"WhenColumn":   whenColumn,
		"Cases":        caseWhens,
		"LastIndex":    len(caseWhens) - 1,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func formatValue(value any) string {
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
