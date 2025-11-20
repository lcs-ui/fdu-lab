package editor

// Editor 编辑器基类（组件接口）

// Command 命令接口（命令模式）
type Command interface {
	Execute()         // 执行命令
	Undo()            // 撤销命令
	IsExecuted() bool // 判断命令是否执行成功
}

// Event 编辑器事件（观察者模式，可选扩展）
type Event struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
	Time int64                  `json:"time"`
}
