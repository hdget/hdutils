package page

import (
	"math"
)

type Utils interface {
	CalculatePages(total, pageSize int64) []Page
	NewPages(total, pageSize int64) []Page
	NewPagination(page, pageSize int64) Pagination
}

type PaginationUtils interface {
	GetLimitClause() string
}

type Page struct {
	Start int64
	End   int64
}

// CalculatePages 根据总数量和每页多少个计算Page数组
func CalculatePages(total, pageSize int64) []Page {
	totalPage := int64(math.Ceil(float64(total) / float64(pageSize)))

	pages := make([]Page, totalPage)
	for i := int64(0); i < totalPage; i++ {
		start := i * pageSize
		end := (i + 1) * pageSize
		if end > total {
			end = total
		}

		pages[i] = Page{
			Start: start,
			End:   end,
		}
	}
	return pages
}

// Paginate 分页
func Paginate[T any](sliceVars []T, pageSize int64) map[int64][]T {
	total := int64(len(sliceVars))
	totalPage := int64(math.Ceil(float64(total) / float64(pageSize)))

	results := make(map[int64][]T)
	for i := int64(0); i < totalPage; i++ {
		start := i * pageSize
		end := (i + 1) * pageSize
		if end > total {
			end = total
		}

		results[i] = sliceVars[start:end]
	}
	return results
}
