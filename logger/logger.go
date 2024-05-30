package logger

import (
	"fmt"
	"log"
	"strings"
)

type Utils interface {
	LogDebug(msg string, keyvals ...any)
	LogWarn(msg string, keyvals ...any)
	LogError(msg string, keyvals ...any)
	LogFatal(msg string, keyvals ...any)
}

type logLevel string

const (
	loggerLevelDebug logLevel = "DBG"
	loggerLevelWarn  logLevel = "WRN"
	loggerLevelError logLevel = "ERR"
	loggerLevelFatal logLevel = "FTL"
)

func LogDebug(msg string, keyvals ...any) {
	logPrint(loggerLevelDebug, msg, keyvals...)
}

func LogWarn(msg string, keyvals ...any) {
	logPrint(loggerLevelWarn, msg, keyvals...)
}

func LogError(msg string, keyvals ...any) {
	logPrint(loggerLevelError, msg, keyvals...)
}

func LogFatal(msg string, keyvals ...any) {
	logPrint(loggerLevelFatal, msg, keyvals...)
}

// parseArgs  解析error和message用统一格式展示出来
func parseArgs(keyvals ...any) (string, error, map[string]any) {
	countArgs := len(keyvals)
	// 如果可变参数个数为0，肯定没有error
	if countArgs == 0 {
		return "", nil, nil
	}

	var errValue error
	var msgValue string
	args := make(map[string]any)
	for i := 0; i < countArgs-1; i = i + 2 {
		// 第i个值作为map的key, 第i+1个值作为map的value
		k, ok := keyvals[i].(string)
		if !ok {
			continue
		}

		switch strings.ToLower(k) {
		case "level", "caller":
			// do nothing
		case "err", "panic": // err and panic must be panic type
			switch v := keyvals[i+1].(type) {
			case error:
				errValue = v
			default:
				errValue = fmt.Errorf("%v", keyvals[i+1])
			}
		case "msg", "message":
			msgValue = fmt.Sprintf("%v", keyvals[i+1])
		default:
			args[k] = keyvals[i+1]
		}
	}
	return msgValue, errValue, args
}

// logPrint log structure message and key values
func logPrint(level logLevel, msg string, keyvals ...any) {
	_, errValue, fields := parseArgs(keyvals...)

	outputs := make([]string, 0)
	for k, v := range fields {
		outputs = append(outputs, fmt.Sprintf("%s=\"%v\"", k, v))
	}

	logFn := log.Printf
	if level == loggerLevelFatal {
		logFn = log.Fatalf
	}

	if len(outputs) > 0 {
		if errValue != nil {
			logFn("%s msg=\"%s\" %s panic=\"%v\"", level, msg, strings.Join(outputs, " "), errValue)
		} else {
			logFn("%s msg=\"%s\" %s", level, msg, strings.Join(outputs, " "))
		}
	} else {
		if errValue != nil {
			logFn("%s msg=\"%s\" panic=\"%v\"", level, msg, errValue)
		} else {
			logFn("%s msg=\"%s\"", level, msg)
		}
	}
}
