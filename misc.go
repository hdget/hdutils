package hdutils

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	"hash/fnv"
	"log"
	"math"
	"os"
	"reflect"
	"runtime/debug"
	"time"
)

// RecordErrorStack 将错误信息保存到错误日志文件中
func RecordErrorStack(app string) {
	filename := fmt.Sprintf("%s.err", app)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer func() {
		if err == nil {
			_ = file.Close()
		}
	}()
	if err != nil {
		log.Printf("RecordErrorStack: error open file, filename: %s, err: %v", filename, err)
		return
	}

	data := bytes.NewBufferString("=== " + time.Now().String() + " ===\n")
	data.Write(debug.Stack())
	data.WriteString("\n")
	_, err = file.Write(data.Bytes())
	if err != nil {
		log.Printf("RecordErrorStack: error write file, filename: %s, err: %v", filename, err)
	}
}

// ReverseInt64Slice 将[]int64 slice倒序重新排列
func ReverseInt64Slice(numbers []int64) []int64 {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

// GetSliceData 将传过来的数据转换成[]interface{}
func GetSliceData(data interface{}) []interface{} {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil
	}

	sliceLenth := v.Len()
	sliceData := make([]interface{}, sliceLenth)
	for i := 0; i < sliceLenth; i++ {
		sliceData[i] = v.Index(i).Interface()
	}

	return sliceData
}

// GetPagePositions 获取分页的起始值列表
// @return 返回一个二维数组， 第一维是多少页，第二维是每页[]int{start, end}
// e,g: 假设11个数的列表，分页pageSize是5，那么返回的是：
//
//	[]int{
//	   []int{0, 5},
//	   []int{5, 10},
//	   []int{10, 11},
//	}
func GetPagePositions(data interface{}, pageSize int) [][]int {
	listData := GetSliceData(data)
	if listData == nil {
		return nil
	}

	total := len(listData)
	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	pages := make([][]int, 0)
	for i := 0; i < totalPage; i++ {
		start := i * pageSize
		end := (i + 1) * pageSize
		if end > total {
			end = total
		}

		p := []int{start, end}
		pages = append(pages, p)
	}
	return pages
}

// GenerateRandString 生成随机字符串
func GenerateRandString(n int) string {
	randStr, _ := gonanoid.Generate(hashids.DefaultAlphabet, n)
	return randStr
}

func HashId(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write(StringToBytes(s))
	return h.Sum32()
}

// GetNeo4jPathPattern 解析Neo4j语法的Variable-length pattern
func GetNeo4jPathPattern(args ...int32) string {
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

func HashString(s string, length int) string {
	return HashBytes(StringToBytes(s), length)
}

func HashBytes(data []byte, length int) string {
	hashValue := fmt.Sprintf("%x", sha256.Sum256(data))
	hdData := hashids.NewData()
	hdData.MinLength = length
	h, _ := hashids.NewWithData(hdData)
	value, _ := h.EncodeHex(hashValue)
	return value
}
