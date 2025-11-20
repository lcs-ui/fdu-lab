package editor

import (
	"strings"
)

// ------------------------------
// 1. 命令接口定义（命令模式核心）
// ------------------------------

// ------------------------------
// 2. AppendCommand：处理 "append" 命令（追加一行）
// ------------------------------

type AppendCommand struct {
	editor    *TextEditor // 关联的编辑器
	text      string      // 要追加的文本（整行）
	prevLines []string    // 追加前的所有行（用于撤销）
	executed  bool        // 是否执行成功
}

// 执行：在文件末尾追加一行

func (cmd *AppendCommand) Execute() {
	if cmd.editor == nil {
		return
	}

	// 保存当前状态（用于撤销）
	cmd.prevLines = make([]string, len(cmd.editor.lines))
	copy(cmd.prevLines, cmd.editor.lines)

	// 执行追加（新增一行）
	cmd.editor.lines = append(cmd.editor.lines, cmd.text)
	cmd.editor.isModified = true
	cmd.executed = true

	// 触发事件（供观察者如日志模块使用）
	//cmd.editor.notifyEvent(Event{
	//	Type: "append",
	//	Data: map[string]interface{}{"text": cmd.text, "line": len(cmd.editor.lines)},
	//	Time: time.Now().UnixMilli(),
	//})

}

// 撤销：删除最后一行（恢复到追加前）

func (cmd *AppendCommand) Undo() {
	if !cmd.executed || cmd.editor == nil {
		return
	}

	// 恢复到追加前的行状态
	cmd.editor.lines = cmd.prevLines
	cmd.editor.isModified = true
}

func (cmd *AppendCommand) IsExecuted() bool {
	return cmd.executed
}

// ------------------------------
// 3. InsertCommand：处理 "insert" 命令（指定位置插入，支持换行）
// ------------------------------

type InsertCommand struct {
	editor     *TextEditor // 关联的编辑器
	line       int         // 目标行号（1-based）
	col        int         // 目标列号（1-based）
	text       string      // 插入的文本（可能含换行符）
	prevLine   string      // 插入前的目标行内容（用于撤销）
	splitLines []string    // 文本按换行拆分后的行（用于执行）
	executed   bool        // 是否执行成功
}

func NewInsertCommand(editor *TextEditor, line, col int, text string) *InsertCommand {
	return &InsertCommand{
		editor: editor,
		line:   line,
		col:    col,
		text:   text,
	}
}

// 执行：在指定位置插入文本（支持换行拆分）

func (cmd *InsertCommand) Execute() {
	if cmd.editor == nil || !cmd.validate() {
		return
	}

	// 转换为 0-based 索引
	lineIdx := cmd.line - 1
	colIdx := cmd.col - 1

	// 保存插入前的行内容（用于撤销）
	cmd.prevLine = cmd.editor.lines[lineIdx]

	// 按换行符拆分文本（支持多行插入）
	cmd.splitLines = strings.Split(cmd.text, "\n")

	// 执行插入逻辑
	if len(cmd.splitLines) == 1 {
		// 无换行：直接插入到当前行
		currentLine := cmd.prevLine
		newLine := currentLine[:colIdx] + cmd.text + currentLine[colIdx:]
		cmd.editor.lines[lineIdx] = newLine
	} else {
		// 有换行：拆分当前行并插入多行
		currentLine := cmd.prevLine
		// 第一部分：当前行从开始到插入位置 + 拆分的第一行
		firstPart := currentLine[:colIdx] + cmd.splitLines[0]
		// 中间部分：拆分的中间行（除首尾外）
		middleParts := cmd.splitLines[1 : len(cmd.splitLines)-1]
		// 最后部分：拆分的最后一行 + 当前行从插入位置到结尾
		lastPart := cmd.splitLines[len(cmd.splitLines)-1] + currentLine[colIdx:]

		// 重组所有行（插入新行）
		newLines := make([]string, 0, len(cmd.editor.lines)+len(middleParts)+1)
		newLines = append(newLines, cmd.editor.lines[:lineIdx]...)   // 插入行之前的内容
		newLines = append(newLines, firstPart)                       // 第一部分
		newLines = append(newLines, middleParts...)                  // 中间部分
		newLines = append(newLines, lastPart)                        // 最后部分
		newLines = append(newLines, cmd.editor.lines[lineIdx+1:]...) // 插入行之后的内容
		cmd.editor.lines = newLines
	}

	cmd.editor.isModified = true
	cmd.executed = true

	// 触发事件
	//cmd.editor.notifyEvent(Event{
	//	Type: "insert",
	//	Data: map[string]interface{}{"line": cmd.line, "col": cmd.col, "text": cmd.text},
	//	Time: time.Now().UnixMilli(),
	//})
}

// 撤销：移除插入的内容（恢复到插入前）

func (cmd *InsertCommand) Undo() {
	if !cmd.executed || cmd.editor == nil {
		return
	}

	lineIdx := cmd.line - 1

	if len(cmd.splitLines) == 1 {
		// 无换行：直接恢复原行
		cmd.editor.lines[lineIdx] = cmd.prevLine
	} else {
		// 有换行：合并被拆分的行，删除插入的中间行
		removeCount := len(cmd.splitLines) - 1 // 需要删除的行数
		newLines := make([]string, 0, len(cmd.editor.lines)-removeCount)
		newLines = append(newLines, cmd.editor.lines[:lineIdx]...)               // 插入行之前的内容
		newLines = append(newLines, cmd.prevLine)                                // 恢复原行
		newLines = append(newLines, cmd.editor.lines[lineIdx+1+removeCount:]...) // 跳过插入的中间行
		cmd.editor.lines = newLines
	}

	cmd.editor.isModified = true
}

// 验证插入位置是否合法
func (cmd *InsertCommand) validate() bool {
	lineCount := len(cmd.editor.lines)

	// 空文件只能在 1:1 位置插入
	if lineCount == 0 {
		return cmd.line == 1 && cmd.col == 1
	}

	// 行号越界（必须在 1~lineCount 之间）
	if cmd.line < 1 || cmd.line > lineCount {
		return false
	}

	// 列号越界（必须在 1~行长度+1 之间，允许插入到行尾）

	targetLine := cmd.editor.lines[cmd.line-1]
	return cmd.col >= 1 && cmd.col <= len(targetLine)+1
}

func (cmd *InsertCommand) IsExecuted() bool {
	return cmd.executed
}

// ------------------------------
// 4. DeleteCommand：处理 "delete" 命令（删除指定长度字符）
// ------------------------------

type DeleteCommand struct {
	editor   *TextEditor // 关联的编辑器
	line     int         // 目标行号（1-based）
	col      int         // 起始列号（1-based）
	length   int         // 删除长度
	prevLine string      // 删除前的行内容（用于撤销）
	executed bool        // 是否执行成功
}

func NewDeleteCommand(editor *TextEditor, line, col, length int) *DeleteCommand {
	return &DeleteCommand{
		editor: editor,
		line:   line,
		col:    col,
		length: length,
	}
}

// 执行：删除指定范围的字符（不可跨行）

func (cmd *DeleteCommand) Execute() {
	if cmd.editor == nil || !cmd.validate() {
		return
	}

	lineIdx := cmd.line - 1
	colIdx := cmd.col - 1

	// 保存删除前的行内容（用于撤销）
	cmd.prevLine = cmd.editor.lines[lineIdx]

	// 执行删除
	currentLine := cmd.prevLine
	newLine := currentLine[:colIdx] + currentLine[colIdx+cmd.length:]
	cmd.editor.lines[lineIdx] = newLine

	cmd.editor.isModified = true
	cmd.executed = true

	// 触发事件
	//cmd.editor.notifyEvent(Event{
	//	Type: "delete",
	//	Data: map[string]interface{}{"line": cmd.line, "col": cmd.col, "length": cmd.length},
	//	Time: time.Now().UnixMilli(),
	//})
}

// 撤销：恢复被删除的字符

func (cmd *DeleteCommand) Undo() {
	if !cmd.executed || cmd.editor == nil {
		return
	}

	// 恢复原行内容
	cmd.editor.lines[cmd.line-1] = cmd.prevLine
	cmd.editor.isModified = true
}

// 验证删除范围是否合法
func (cmd *DeleteCommand) validate() bool {
	lineCount := len(cmd.editor.lines)

	// 行号越界
	if cmd.line < 1 || cmd.line > lineCount {
		return false
	}

	targetLine := cmd.editor.lines[cmd.line-1]
	lineLen := len(targetLine)
	colIdx := cmd.col - 1

	// 列号越界或删除长度无效
	if colIdx < 0 || colIdx >= lineLen || cmd.length <= 0 {
		return false
	}

	// 删除范围不能超过行尾
	if colIdx+cmd.length > lineLen {
		return false
	}

	return true
}

func (cmd *DeleteCommand) IsExecuted() bool {
	return cmd.executed
}

// ------------------------------
// 5. ReplaceCommand：处理 "replace" 命令（先删后插）
// ------------------------------

type ReplaceCommand struct {
	editor    *TextEditor    // 关联的编辑器
	line      int            // 目标行号（1-based）
	col       int            // 起始列号（1-based）
	length    int            // 删除长度
	text      string         // 替换的新文本
	deleteCmd *DeleteCommand // 内部删除命令
	insertCmd *InsertCommand // 内部插入命令
	executed  bool           // 是否执行成功
}

func NewReplaceCommand(editor *TextEditor, line, col, length int, text string) *ReplaceCommand {
	return &ReplaceCommand{
		editor:    editor,
		line:      line,
		col:       col,
		length:    length,
		text:      text,
		deleteCmd: NewDeleteCommand(editor, line, col, length),
		insertCmd: NewInsertCommand(editor, line, col, text), // 插入位置与删除位置相同
	}
}

// 执行：先删除指定长度字符，再插入新文本

func (cmd *ReplaceCommand) Execute() {
	if cmd.editor == nil {
		return
	}

	// 先执行删除
	cmd.deleteCmd.Execute()
	if !cmd.deleteCmd.IsExecuted() {
		return // 删除失败则终止替换
	}

	// 再执行插入（删除后行结构可能变化，但插入位置仍基于原行号）
	cmd.insertCmd.Execute()
	cmd.executed = cmd.insertCmd.IsExecuted()

	//// 触发事件
	//cmd.editor.notifyEvent(Event{
	//	Type: "replace",
	//	Data: map[string]interface{}{"line": cmd.line, "col": cmd.col, "length": cmd.length, "text": cmd.text},
	//	Time: time.Now().UnixMilli(),
	//})
}

// 撤销：先撤销插入，再撤销删除（恢复原状态）

func (cmd *ReplaceCommand) Undo() {
	if !cmd.executed || cmd.editor == nil {
		return
	}

	// 先撤销插入（移除新文本）
	cmd.insertCmd.Undo()
	// 再撤销删除（恢复原文本）
	cmd.deleteCmd.Undo()

	cmd.editor.isModified = true
}

func (cmd *ReplaceCommand) IsExecuted() bool {
	return cmd.executed
}
