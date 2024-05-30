package neo4j

import "fmt"

// GetPathPattern 解析Neo4j语法的Variable-length pattern
func GetPathPattern(args ...int32) string {
	start := int32(-1)
	end := int32(-1)
	switch len(args) {
	case 1:
		start = args[0]
	case 2:
		start = args[0]
		end = args[1]
	}

	expr := "*"
	if start >= 0 {
		if end >= start {
			expr = fmt.Sprintf("*%d..%d", start, end)
		} else {
			expr = fmt.Sprintf("*%d..", start)
		}
	}
	return expr
}
