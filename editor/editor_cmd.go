package editor

import (
	"fmt"
	"lab1/common"
	"strings"
	"time"
	"strconv"
)

func NewAppendCommand(editor *TextEditor, text string) *AppendCommand {
	return &AppendCommand{
		editor: editor,
		text:   text,
	}
}

// 暴露给外部的操作方法（供用户指令调用）

func (te *TextEditor) Append(text string) {
	if te.logEnabled{
	te.workspaceApi.NotifyObservers(common.WorkspaceEvent{
		FilePath: te.GetFilePath(),
		Type: "Append",
		Command: "Append " + text,
		Timestamp: time.Now().UnixMilli(),
	})		
	}

	te.ExecuteCommand(NewAppendCommand(te, text))
}

func (te *TextEditor) Insert(line, col int, text string) {
	if te.logEnabled{
	commandStr := "Insert " + strconv.Itoa(line) + "," + strconv.Itoa(col) + " " + text
	te.workspaceApi.NotifyObservers(common.WorkspaceEvent{
		FilePath: te.GetFilePath(),
		Type:     "Insert",
		Command:  commandStr,
		Timestamp: time.Now().UnixMilli(),
	})
	}
	te.ExecuteCommand(NewInsertCommand(te, line, col, text))
}

func (te *TextEditor) Delete(line, col, length int) {
	if te.logEnabled{
	commandStr := "Delete "+ strconv.Itoa(line)+","+strconv.Itoa(col)+","+strconv.Itoa(length)
	te.workspaceApi.NotifyObservers(common.WorkspaceEvent{
		FilePath: te.GetFilePath(),
		Type: "Delete",
		Command: commandStr,
		Timestamp: time.Now().UnixMilli(),
	})
	}
	te.ExecuteCommand(NewDeleteCommand(te, line, col, length))
}

func (te *TextEditor) Replace(line, col, length int, text string) {
	if te.logEnabled{
	commandStr := "Replace "+ strconv.Itoa(line)+","+strconv.Itoa(col)+","+strconv.Itoa(length)+" "+text
	te.workspaceApi.NotifyObservers(common.WorkspaceEvent{
		FilePath: te.GetFilePath(),
		Type: "Relpace",
		Command: commandStr,
		Timestamp: time.Now().UnixMilli(),
	})
	}
	te.ExecuteCommand(NewReplaceCommand(te, line, col, length, text))
}

// Show 方法
func (te *TextEditor) Show(startLine, endLine int) {
	if te.logEnabled{
	commandStr := "Show "+ strconv.Itoa(startLine)+","+strconv.Itoa(endLine)
	te.workspaceApi.NotifyObservers(common.WorkspaceEvent{
		FilePath: te.GetFilePath(),
		Type: "Show",
		Command: commandStr,
		Timestamp: time.Now().UnixMilli(),
	})
	}

	lineCount := len(te.lines)

	// 处理空文件
	if lineCount == 0 {
		fmt.Println("(空文件)")
		return
	}

	// 解析行范围（默认显示全文）
	actualStart := 1
	actualEnd := lineCount

	if startLine > 0 {
		// 修正起始行越界（最小为1）
		actualStart = startLine
		if actualStart < 1 {
			actualStart = 1
		}
		// 修正起始行超过总行数（视为无效范围）
		if actualStart > lineCount {
			fmt.Println("起始行超出文件范围")
			return
		}

		// 修正结束行（默认到最后一行，最大为总行数）
		if endLine > 0 {
			actualEnd = endLine
			if actualEnd > lineCount {
				actualEnd = lineCount
			}
			// 起始行不能大于结束行
			if actualStart > actualEnd {
				fmt.Println("起始行不能大于结束行")
				return
			}
		}
	}

	// 计算行号宽度（用于对齐）
	maxLineNum := actualEnd
	if lineCount > maxLineNum {
		maxLineNum = lineCount // 确保行号宽度适配最大行号
	}
	lineNumWidth := len(fmt.Sprintf("%d", maxLineNum))
	lineFormat := fmt.Sprintf("%%%dd: %%s\n", lineNumWidth) // 格式：  1: Hello

	//
	// 拼接输出内容
	var output strings.Builder
	for i := actualStart - 1; i < actualEnd; i++ { // 转换为 0-based 索引
		lineNum := i + 1
		output.WriteString(fmt.Sprintf(lineFormat, lineNum, te.lines[i]))
	}

	// 打印结果（去除末尾多余换行）
	fmt.Print(output.String())
}
