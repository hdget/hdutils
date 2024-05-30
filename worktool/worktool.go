package worktool

import (
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
)

type Utils interface {
	SendFile(targets []string, file File) error
	SendText(targets []string, content string, atPeoples ...string) error
}

type apiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type workToolImpl struct {
	robotId string // 机器人id
}

const (
	socketType = 2 //通讯类型
	apiUrl     = "https://worktool.asrtts.cn/wework/sendRawMessage"
)

func New(robotId string) Utils {
	return &workToolImpl{robotId: robotId}
}

func (w *workToolImpl) SendText(targets []string, content string, atPeoples ...string) error {
	msg := &textMessage{
		Kind:      msgKindText,
		Targets:   targets,
		Content:   content,
		AtPeoples: atPeoples,
	}

	return w.send(msg)
}

func (w *workToolImpl) SendFile(targets []string, file File) error {
	msg := fileMessage{
		Kind:    msgKindFile,
		Targets: targets,
		File:    file,
	}

	return w.send(msg)
}

func (w *workToolImpl) send(messages ...any) error {
	body := messageBody{
		SocketType: socketType,
		List:       messages,
	}

	var resp apiResponse
	ret, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("robotId", w.robotId).
		SetBody(body).
		SetResult(&resp).
		Post(apiUrl)
	if err != nil {
		return err
	}

	// http status code != 200
	if ret.StatusCode() != 200 {
		return errors.New(convert.BytesToString(ret.Body()))
	}

	// api response code != 200 means error
	if resp.Code != 200 {
		return errors.New(resp.Message)
	}

	return nil
}
