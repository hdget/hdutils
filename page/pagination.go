package page

import (
	"fmt"
)

const (
	DefaultPageSize = 10
)

type Pagination struct {
	Page     uint64 // 第几页
	PageSize uint64 // 每页几项
	Offset   uint64 // 偏移起始值
}

// NewPagination 获取分页对象
func NewPagination(page, pageSize int64) Pagination {
	// 处理当前页面
	if page <= 0 {
		page = 1
	}

	// 处理每页大小
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	offset := uint64((page - 1) * pageSize)

	return Pagination{Page: uint64(page), PageSize: uint64(pageSize), Offset: offset}
}

// GetStartEndPosition 如果是按列表slice进行翻页的话， 计算slice的起始index
func GetStartEndPosition(page, pageSize, total int64) (int64, int64) {
	// 处理当前页面
	if page <= 0 {
		page = 1
	}

	// 处理每页大小
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	start := (page - 1) * pageSize
	end := page * pageSize

	if end > total {
		end = total
	}

	if start > end {
		start = end
	}

	return start, end
}

// GetLimitClause 获取limit sql子句
func (p Pagination) GetLimitClause() string {
	if p.Offset == 0 {
		return fmt.Sprintf("LIMIT %d", p.PageSize)
	}

	return fmt.Sprintf("LIMIT %d, %d", p.Offset, p.PageSize)
}

// GetSQLClause 获取翻页SQL查询语句
//
// 1. 假如前端没有传过来last_pk, 那么返回值是 last_pk, LIMIT子句(LIMIT offset, PageSize)
// e,g: 0, "LIMIT 20, 10" => 在数据库查询时可能会被组装成 WHERE pk > 0 ...  LIMIT 20, 10
//
// 2. 假如前端传过来了last_pk, 那么返回值是 last_pk, LIMIT子句(LIMIT PageSize)
// e,g: 123,"LIMIT 10" => 在数据库查询时可能会被组装成 WHERE pk > 123 ...  LIMIT 10
//func (p *Pagination) GetSQLClause(total int64) text {
//	if p == nil {
//		return ""
//	}
//
//	// 如果total值为0, 默认返回指定页面
//	if total == 0 {
//		return "LIMIT 0"
//	}
//
//	start := (p.Page - 1) * p.PageSize
//	return fmt.Sprintf("LIMIT %d, %d", start, p.PageSize)
//	//start, end := GetStartEndPosition(p.Page, p.PageSize, total)
//	//
//	//return fmt.Sprintf("LIMIT %d, %d", start, end-start)
//}
// GetPagePositions 获取分页的起始值列表
// @return 返回一个二维数组， 第一维是多少页，第二维是每页[]int{start, end}
// e,g: 假设11个数的列表，分页pageSize是5，那么返回的是：
//
//	[]int{
//	   []int{0, 5},
//	   []int{5, 10},
//	   []int{10, 11},
//	}

//func GetPagePositions(data any, pageSize int) [][]int {
//	listData := GetSliceData(data)
//	if listData == nil {
//		return nil
//	}
//
//	total := len(listData)
//	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))
//
//	pages := make([][]int, 0)
//	for i := 0; i < totalPage; i++ {
//		start := i * pageSize
//		end := (i + 1) * pageSize
//		if end > total {
//			end = total
//		}
//
//		p := []int{start, end}
//		pages = append(pages, p)
//	}
//	return pages
//}
