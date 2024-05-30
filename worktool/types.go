package worktool

type File struct {
	Name    string `json:"objectName"` // 对象名，这里代表文件名
	Url     string `json:"fileUrl"`    // 文件路径
	Kind    string `json:"fileType"`   // 文件类型 image等
	Comment string `json:"extraText"`  // 附件流言 可选填
}

type textMessage struct {
	Kind      msgKind  `json:"type"`            // 消息类型
	Targets   []string `json:"titleList"`       // 发送目标，可以为联系人或者群
	Content   string   `json:"receivedContent"` // 文本内容 (\n换行)
	AtPeoples []string `json:"atList"`          // @人列表
}

type fileMessage struct {
	File
	Kind    msgKind  `json:"type"`
	Targets []string `json:"titleList"` // 发送人列表
}

type messageBody struct {
	SocketType int   `json:"socketType"`
	List       []any `json:"list"`
}

type msgKind int

const (
	msgKindText msgKind = 203
	msgKindFile msgKind = 218
)
