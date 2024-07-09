package sql

import (
	"fmt"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
	"strings"
)

type BatchUpdater interface {
	Set(setColumn string, value any) BatchUpdater
	Case(caseSetColumn, caseWhenColumn string) BatchUpdater
	When(whenValue, thenValue any) BatchUpdater
	Generate() (string, error)
}

const (
	templateBatchUpdate = `
UPDATE %s SET %s = CASE 
%s 
END WHERE %s IN (%s)`
)

type mysqlBatchUpdater struct {
	table          string
	caseSetColumn  string        // case when中更新的字段
	caseWhenColumn string        // case when中比较的字段
	otherSets      []*setClause  // 其他同时更新的子句
	whens          []*whenClause // 条件更新的子句
}

type setClause struct {
	column string
	value  any
}

type whenClause struct {
	whenValue any
	thenValue any
}

func NewMysqlBatchUpdater(table string) BatchUpdater {
	return &mysqlBatchUpdater{
		table:     table,
		otherSets: make([]*setClause, 0),
		whens:     make([]*whenClause, 0),
	}
}

// Set 批量更新的时候同时更新的字段
func (u *mysqlBatchUpdater) Set(setColumn string, value any) BatchUpdater {
	u.otherSets = append(u.otherSets, &setClause{
		column: setColumn,
		value:  value,
	})
	return u
}

// Case 更新和比较的字段
func (u *mysqlBatchUpdater) Case(caseSetColumn, caseWhenColumn string) BatchUpdater {
	u.caseSetColumn = caseSetColumn
	u.caseWhenColumn = caseWhenColumn
	return u
}

// When 条件更新
func (u *mysqlBatchUpdater) When(whenValue, thenValue any) BatchUpdater {
	u.whens = append(u.whens, &whenClause{
		whenValue: whenValue,
		thenValue: thenValue,
	})
	return u
}

func (u *mysqlBatchUpdater) Generate() (string, error) {
	if u.table == "" || u.caseSetColumn == "" || len(u.whens) == 0 {
		return "", errors.New("invalid parameter")
	}

	setParts := make([]string, 0)
	// 如果有额外要更新的字段，组装起来
	for _, otherSet := range u.otherSets {
		setParts = append(setParts, fmt.Sprintf("%s=%s", u.escape(otherSet.column), u.formatValue(otherSet.value)))
	}

	// 添加条件更新的字段
	setParts = append(setParts, u.escape(u.caseSetColumn))

	// 构造when部分
	whenParts := make([]string, 0)
	whereInParts := make([]string, 0)
	for _, when := range u.whens {
		whenParts = append(whenParts, fmt.Sprintf("\tWHEN %s = %s THEN %s", u.escape(u.caseWhenColumn), u.formatValue(when.whenValue), u.formatValue(when.thenValue)))
		whereInParts = append(whereInParts, u.formatValue(when.whenValue))
	}

	return fmt.Sprintf(
		templateBatchUpdate,
		u.escape(u.table),
		strings.Join(setParts, ","),    // 填入set部分
		strings.Join(whenParts, " \n"), // 填入when部分
		u.escape(u.caseWhenColumn),
		strings.Join(whereInParts, ","), // 填入where in部分
	), nil
}

func (u *mysqlBatchUpdater) escape(s string) string {
	return fmt.Sprintf("`%s`", s)
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
