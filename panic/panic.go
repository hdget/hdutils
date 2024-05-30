package panic

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"
)

// RecordErrorStack 将错误信息保存到错误日志文件中
func RecordErrorStack(app string) {
	filename := fmt.Sprintf("%s.dump", app)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer func() {
		if err == nil {
			_ = file.Close()
		}
	}()
	if err != nil {
		log.Printf("RecordErrorStack: panic open file, filename: %s, err: %v", filename, err)
		return
	}

	data := bytes.NewBufferString("=== " + time.Now().String() + " ===\n")
	data.Write(debug.Stack())
	data.WriteString("\n")
	_, err = file.Write(data.Bytes())
	if err != nil {
		log.Printf("RecordErrorStack: panic write file, filename: %s, err: %v", filename, err)
	}
}
