package main

import (
	"bufio"
	"fmt"
	"lab1/common"
	"lab1/editor"
	"lab1/log"
	"lab1/storage"
	"lab1/treeview"
	"lab1/workspace"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// 1. 初始化依赖组件
	fileStorage := storage.NewLocalStorage("./workspace_state.json") // 状态存储路径
	logModule := log.NewLogModule()

	// 2. 初始化工作区
	ws := workspace.NewWorkspace("./workspace_state.json")

	// 3. 日志模块订阅工作区事件（观察者模式）
	ws.RegisterObserver(logModule)

	//日志模块订阅编辑器事件

	// 4. 从本地存储恢复上次工作区状态（备忘录模式）
	if err := restoreWorkspaceState(ws, fileStorage); err != nil {
		fmt.Printf("恢复工作区失败，使用新状态: %v\n", err)
	} else {
		fmt.Println("工作区已恢复上次状态")
	}

	// 5. 启动交互循环，处理用户指令
	startInteractiveLoop(ws)
}

// 修复后的 restoreWorkspaceState 函数
func restoreWorkspaceState(ws *workspace.Workspace, storage *storage.LocalStorage) error {
	// 调用 Workspace 的 RestoreState 方法，传入编辑器工厂函数
	// 工厂函数复用之前定义的 editor.EditorFactory（需确保已导入 editor 包）
	return ws.RestoreState(editor.EditorFactory)
}

// 启动用户交互循环
func startInteractiveLoop(ws *workspace.Workspace) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("编辑器启动完成，支持指令: load/save/close/undo/exit....")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		handleCommand(ws, input, true)
		//fmt.Printf("[debug]active_file: %s\n", ws.GetActiveEditor().GetFilePath())
		activeEditor := ws.GetActiveEditor()
		if activeEditor == nil {
			fmt.Println("[debug]active_file: 无激活的编辑器/文件")
		} else {
			fmt.Printf("[debug]active_file: %s\n", activeEditor.GetFilePath())
		}
	}
}

// 处理用户指令
func handleCommand(ws *workspace.Workspace, input string, debug bool) {
	parts := strings.SplitN(input, " ", 4)
	if len(parts) == 0 {
		fmt.Println("无效指令")
		return
	}
	cmd := parts[0]
	switch cmd {
	case "load": //完成
		_load(ws, parts, false)
	case "save": //完成
		_Save(ws, input, debug, parts)
	case "close": //完成
		_close(ws, parts)
	case "init": //完成
		_init(ws, parts)
	case "undo":
		_undo(ws)
	case "redo":
		_redo(ws)
	case "editor-list": //完成
		_EditorList(ws)
	case "edit": //完成
		_edit(ws, parts)
	case "exit":
		_exit(ws)
	case "dir-tree": //完成
		_dirTree(ws, parts)
	case "xml-tree": //完成
		_xmlTree(ws, parts)
	case "append":
		_append(ws, parts)
	case "insert":
		_insert(ws, parts)
	case "show":
		_show(ws, parts)
	case "delete":
		_delete(ws, parts)
	case "replace":
		_replace(ws, parts)
	case "log-on":
		_LogOn(ws, parts)
	case "log-off":
		_LogOff(ws, parts)
	case "log-show":
		_LogShow(ws, parts)
	default:
		fmt.Println("未知指令，支持: load/save/close/undo/exit")
	}
}
func _load(ws *workspace.Workspace, parts []string, debug bool) {
	if len(parts) < 2 {
		fmt.Println("请指定文件路径: load [path]")
		return
	}
	_editor, err := ws.LoadFile(parts[1], editor.EditorFactory)
	if err != nil {
		fmt.Printf("加载失败: %v\n", err)
	} else {
		fmt.Printf("已加载文件: %s（%s）\n",
			_editor.GetFilePath(),
			map[bool]string{true: "已修改", false: "未修改"}[_editor.IsModified()])
		if debug {
			fmt.Println("[debug] 当前活动文件是" + ws.GetActiveEditor().GetFilePath())
		}
	}
}

func _close(ws *workspace.Workspace, parts []string) {
	if len(parts) < 2 {
		fmt.Println("请指定文件路径: close [path]")
		return
	}
	if err := ws.CloseFile(parts[1]); err != nil {
		fmt.Printf("关闭失败: %v\n", err)
	} else {
		fmt.Printf("已关闭文件: %s\n", parts[1])
	}
}

func _undo(ws *workspace.Workspace) {
	if err := ws.GetActiveEditor().Undo(); err != nil {
		fmt.Printf("undo失败: %v\n", err)
	} else {
		fmt.Println("undo成功")
	}

}

func _redo(ws *workspace.Workspace) {
	if err := ws.GetActiveEditor().Redo(); err != nil {
		fmt.Printf("redo失败: %v\n", err)
	} else {
		fmt.Println("redo成功")
	}
}

func _exit(ws *workspace.Workspace) {
	// 退出前保存工作区状态
	memento := ws.CreateMemento()
	if err := storage.NewLocalStorage("./workspace_state.json").SaveMemento(memento); err != nil {
		fmt.Printf("保存工作区状态失败: %v\n", err)
	}
	fmt.Println("程序退出")
	os.Exit(0)
}

func _dirTree(ws *workspace.Workspace, parts []string) {
	// 确定目标目录（默认当前工作目录）
	targetDir := "."
	if len(parts) >= 2 {
		targetDir = parts[1]
	}

	// 验证目录是否存在
	if _, err := os.Stat(targetDir); err != nil {
		fmt.Printf("目录不存在: %v\n", err)
		return
	}

	// 使用适配器模式生成目录树
	adapter := treeview.NewDirectoryTreeAdapter(targetDir)
	rootNode, err := adapter.GetRootNode()
	if err != nil {
		fmt.Printf("生成目录树失败: %v\n", err)
		return
	}
	
	// 使用 TreeView 渲染树形结构
	view := treeview.NewTreeView(rootNode)
	fmt.Print(view.Render())
}

func _xmlTree(ws *workspace.Workspace, parts []string) {
	// 确定目标XML文件
	if len(parts) < 2 {
		fmt.Println("请指定XML文件路径: xml-tree [path]")
		return
	}
	
	xmlPath := parts[1]
	
	// 如果路径不是绝对路径，尝试从 ./files 目录读取
	if !filepath.IsAbs(xmlPath) && !strings.HasPrefix(xmlPath, ".") {
		xmlPath = filepath.Join("./files", xmlPath)
	}
	
	// 验证文件是否存在
	if _, err := os.Stat(xmlPath); err != nil {
		fmt.Printf("文件不存在: %v\n", err)
		return
	}
	
	// 使用适配器模式生成XML树
	adapter := treeview.NewXMLTreeAdapter(xmlPath)
	rootNode, err := adapter.GetRootNode()
	if err != nil {
		fmt.Printf("解析XML文件失败: %v\n", err)
		return
	}
	
	// 使用 TreeView 渲染树形结构
	view := treeview.NewTreeView(rootNode)
	fmt.Print(view.Render())
}

func _LogOn(ws *workspace.Workspace, parts []string) {
	targetEditor := getTargetEditor(ws, parts) // 解析目标文件（见下方辅助函数）
	if targetEditor == nil {
		fmt.Println("错误：文件未找到或无活动文件")
		return
	}
	targetEditor.SetLogEnabled(true)
	fmt.Printf("已为文件 %s 启用日志\n", targetEditor.GetFilePath())
}

// 处理log-off：关闭指定文件/当前活动文件的日志
func _LogOff(ws *workspace.Workspace, parts []string) {
	targetEditor := getTargetEditor(ws, parts)
	if targetEditor == nil {
		fmt.Println("错误：文件未找到或无活动文件")
		return
	}
	targetEditor.SetLogEnabled(false)
	fmt.Printf("已关闭文件 %s 的日志\n", targetEditor.GetFilePath())
}

// 处理log-show：显示指定文件/当前活动文件的日志
func _LogShow(ws *workspace.Workspace, parts []string) {
	targetEditor := getTargetEditor(ws, parts)
	if targetEditor == nil {
		fmt.Println("错误：文件未找到或无活动文件")
		return
	}

	//fmt.Printf("s%",logFilePath)

	// 打印原始文件路径和计算的日志路径（用于调试）
	fmt.Printf("调试：目标文件路径 = %q\n", targetEditor.GetFilePath())
	// logFilePath := "." + filePath + ".log"
	logFilePath := "." + targetEditor.GetFilePath() + ".log"
	fmt.Printf("调试：日志文件路径 = %q\n", logFilePath) // 检查路径是否正确

	content, err := os.ReadFile(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("日志文件不存在：%s\n", logFilePath)
			return
		}
		fmt.Printf("读取日志失败：%v\n", err)
		return
	}
	fmt.Printf("===== 日志内容（%s） =====\n", logFilePath)
	fmt.Print(string(content))
}

// 辅助函数：获取目标文件的编辑器（支持指定文件或当前活动文件）
func getTargetEditor(ws *workspace.Workspace, parts []string) common.Editor {
	if len(parts) >= 2 {
		// 指定文件：从已打开的编辑器中查找
		if editor, exists := ws.OpenEditors[parts[1]]; exists {
			return editor
		}
		return nil
	} else {
		// 无参数：使用当前活动文件
		return ws.GetActiveEditor()
	}
}

func _Save(ws *workspace.Workspace, input string, debug bool, parts []string) {
	if debug {
		fmt.Printf("[DEBUG] 进入 save 命令处理，输入: %q，参数拆分: %v\n", input, parts)
	}

	// 1. 处理无参数：保存当前活动文件
	if len(parts) == 1 {
		if debug {
			fmt.Println("[DEBUG] 无参数，尝试保存当前活动文件")
		}
		activeEditor := ws.GetActiveEditor()
		if activeEditor == nil {
			if debug {
				fmt.Println("[DEBUG] 未找到活动文件")
			}
			fmt.Println("没有活动文件可保存")
			return
		}
		if debug {
			fmt.Printf("[DEBUG] 找到活动文件: %s，准备保存\n", activeEditor.GetFilePath())
		}
		if err := ws.SaveFile(activeEditor); err != nil {
			if debug {
				fmt.Printf("[DEBUG] 活动文件保存失败: %v\n", err)
			}
			fmt.Printf("保存失败: %v\n", err)
		} else {
			if debug {
				fmt.Printf("[DEBUG] 活动文件保存成功: %s\n", activeEditor.GetFilePath())
			}
			fmt.Printf("已保存活动文件: %s\n", activeEditor.GetFilePath())
		}
		return
	}

	// 2. 处理参数：保存指定文件或所有文件
	subCmd := parts[1]
	if debug {
		fmt.Printf("[DEBUG] 检测到子命令: %q\n", subCmd)
	}
	switch subCmd {
	case "all":
		if debug {
			fmt.Println("[DEBUG] 处理 save all，准备保存所有打开的文件")
		}
		// 保存所有已打开的文件
		openEditors := ws.GetOpenEditors()
		if len(openEditors) == 0 {
			if debug {
				fmt.Println("[DEBUG] 未找到任何打开的文件")
			}
			fmt.Println("没有打开的文件可保存")
			return
		}
		if debug {
			fmt.Printf("[DEBUG] 共找到 %d 个打开的文件，开始批量保存\n", len(openEditors))
		}
		successCount := 0
		for i, editor := range openEditors {
			if debug {
				fmt.Printf("[DEBUG] 正在保存第 %d 个文件: %s\n", i+1, editor.GetFilePath())
			}
			if err := ws.SaveFile(editor); err != nil {
				if debug {
					fmt.Printf("[DEBUG] 第 %d 个文件保存失败: %v\n", i+1, err)
				}
				fmt.Printf("保存文件 %s 失败: %v\n", editor.GetFilePath(), err)
			} else {
				successCount++
				if debug {
					fmt.Printf("[DEBUG] 第 %d 个文件保存成功: %s\n", i+1, editor.GetFilePath())
				}
			}
		}
		if debug {
			fmt.Printf("[DEBUG] 批量保存完成，成功 %d 个，失败 %d 个\n", successCount, len(openEditors)-successCount)
		}
		fmt.Printf("批量保存完成，成功 %d 个，失败 %d 个\n", successCount, len(openEditors)-successCount)

	default:
		// 保存指定文件（subCmd 为文件路径）
		targetPath := subCmd
		if debug {
			fmt.Printf("[DEBUG] 处理指定文件保存，目标路径: %s\n", targetPath)
		}
		// 检查文件是否已打开
		openEditors := ws.GetOpenEditors()
		var targetEditor common.Editor
		for _, editor := range openEditors {
			if editor.GetFilePath() == targetPath {
				targetEditor = editor
				if debug {
					fmt.Printf("[DEBUG] 在已打开文件中找到目标文件: %s\n", targetPath)
				}
				break
			}
		}
		if targetEditor == nil {
			if debug {
				fmt.Printf("[DEBUG] 目标文件 %s 未打开\n", targetPath)
			}
			fmt.Printf("文件 %s 未打开，无法保存\n", targetPath)
			return
		}
		// 执行保存
		if err := ws.SaveFile(targetEditor); err != nil {
			if debug {
				fmt.Printf("[DEBUG] 指定文件 %s 保存失败: %v\n", targetPath, err)
			}
			fmt.Printf("保存文件 %s 失败: %v\n", targetPath, err)
		} else {
			if debug {
				fmt.Printf("[DEBUG] 指定文件 %s 保存成功\n", targetPath)
			}
			fmt.Printf("已保存文件: %s\n", targetPath)
		}
	}
}

func _init(ws *workspace.Workspace, parts []string) {
	if len(parts) < 2 {
		fmt.Println("用法: init <file> [with-log]")
		return
	}
	fileName := parts[1]
	withLog := len(parts) >= 3 && parts[2] == "with-log"

	// 初始化文件内容
	content := ""
	if withLog {
		content = "# log\n" // 带日志标记的初始化内容
	}
	fmt.Printf("[%s %s %s]", parts[0], parts[1], parts[2])
	// 创建未保存的缓冲区（使用 fileName 作为唯一标识）
	_editor := editor.NewTextEditor(fileName, content, ws)
	_editor.MarkAsModified(true) // 新缓冲区默认标记为已修改

	// 添加到工作区的未保存缓冲区，并设为活动文件
	ws.OpenEditors[_editor.GetFilePath()] = _editor
	ws.SetActiveEditor(_editor)

	fmt.Printf("已创建新缓冲区: %s（未保存）\n", fileName)
	if withLog {
		fmt.Println("已自动添加日志标记 '# log'")
	}
}

func _EditorList(ws *workspace.Workspace) {
	openEditors := ws.GetOpenEditors()
	if len(openEditors) == 0 {
		fmt.Printf("error:")
	}

	//j := len(openEditors)
	//fmt.Printf("\n")
	//modified := false //是否修改
	for _, _editor := range openEditors {
		if _editor.GetFilePath() != "" {
			if _editor.IsModified() {
				fmt.Printf("%s [modified]\n", _editor.GetFilePath())
			} else {
				fmt.Printf("%s\n", _editor.GetFilePath())
			}
		}
	}

}

func _edit(ws *workspace.Workspace, parts []string) {
	if len(parts) < 2 {
		fmt.Printf("请指定文件:edit [file]\n")
		return
	}
	fileName := parts[1]
	if fileName == "" {
		fmt.Printf("请指定文件:edit [file]\n")
	} else {
		_, exists := ws.OpenEditors[fileName]
		if exists {
			ws.SetActiveEditor(ws.OpenEditors[fileName])
		} else {
			fmt.Printf("文件未打开: [file]\n")
		}
	}

}

func _show(ws *workspace.Workspace, parts []string) {
	activeEditor := ws.GetActiveEditor()
	if activeEditor == nil {
		fmt.Println("没有活动文件")
		return
	}
	if len(parts) == 1 {
		fmt.Printf("指令格式错误:show [startLine:endLine]\n")
		return
	}
	startLine, endLine := 0, 0
	if len(parts) > 0 {
		rangeStr := parts[1]
		// 按 ":" 分割字符串，处理 "start:end" 格式
		segments := strings.Split(rangeStr, ":")
		if len(segments) != 2 {
			fmt.Println("参数格式错误，应为 show [startLine:endLine]")
			return
		}

		// 解析起始行（必须为正整数）
		s, err := strconv.Atoi(segments[0])
		if err != nil || s < 1 {
			fmt.Println("起始行必须为正整数")
			return
		}

		// 解析结束行（必须为正整数且不小于起始行）
		e, err := strconv.Atoi(segments[1])
		if err != nil || e < 1 {
			fmt.Println("结束行必须为正整数")
			return
		}
		if e < s {
			fmt.Println("结束行不能小于起始行")
			return
		}

		startLine, endLine = s, e
		// 调用编辑器的 Show 方法
		activeEditor.Show(startLine, endLine)
	}

	// 调用编辑器的 Show 方法
	//activeEditor.Show(startLine, endLine)
}

func _append(ws *workspace.Workspace, parts []string) {
	// 1. 校验活动文件是否存在
	activeEditor := ws.GetActiveEditor()
	if activeEditor == nil {
		fmt.Println("错误：没有打开的文件，请先使用 load 命令加载文件")
		return
	}

	// 2. 解析参数：实际参数是 parts[1:]（排除 parts[0] 的 "append"）
	// 检查是否提供了参数（至少需要一个参数片段）
	if len(parts) < 2 { // parts 长度至少为 2（["append", "参数"]）
		fmt.Println("参数错误：请指定要追加的文本，格式为 append \"text\"")
		return
	}

	// 提取实际参数部分（排除命令名），并合并成完整字符串
	// 例如 parts[1:] 是 ["\"hello", "world\""] → 合并为 "\"hello world\""
	textArg := strings.Join(parts[1:], " ")

	// 3. 校验文本是否用双引号包裹
	if len(textArg) < 2 || textArg[0] != '"' || textArg[len(textArg)-1] != '"' {
		fmt.Println("参数错误：文本必须用双引号包裹，格式为 append \"text\"")
		return
	}

	// 4. 提取引号内的文本（去除首尾引号）
	content := textArg[1 : len(textArg)-1]

	// 5. 执行追加操作
	activeEditor.Append(content)
	fmt.Printf("已在文件末尾追加一行：%s\n", content)

}

func _insert(ws *workspace.Workspace, parts []string) {
	// 1. 校验活动文件是否存在
	activeEditor := ws.GetActiveEditor()
	if activeEditor == nil {
		fmt.Println("错误：没有打开的文件，请先使用 load 命令加载文件")
		return
	}

	// 2. 校验参数数量    // 格式要求：至少需要两个参数（位置 <line:col> 和带引号的文本）
	if len(parts) < 3 {
		fmt.Println("参数错误：格式为 insert <line:col> \"text\"（例如 insert 1:4 \"XYZ\"）")
		return
	}

	// 3. 解析位置参数 <line:col>
	posStr := parts[1]
	var line, col int
	// 按 ":" 分割行号和列号
	posParts := strings.Split(posStr, ":")
	if len(posParts) != 2 {
		fmt.Println("参数错误：位置格式应为 line:col（例如 1:4）")
		return
	}
	// 转换行号为整数（1-based）
	line, err := strconv.Atoi(posParts[0])
	if err != nil || line < 1 {
		fmt.Println("参数错误：行号必须为正整数")
		return
	}
	// 转换列号为整数（1-based）
	col, err = strconv.Atoi(posParts[1])
	if err != nil || col < 1 {
		fmt.Println("参数错误：列号必须为正整数")
		return
	}

	// 4. 解析插入文本（合并后续参数，支持带空格和换行符）
	// 文本参数从 parts[2] 开始（排除命令名和位置参数）
	textParts := parts[2:]
	textArg := strings.Join(textParts, " ")
	// 校验文本是否用双引号包裹
	if len(textArg) < 2 || textArg[0] != '"' || textArg[len(textArg)-1] != '"' {
		fmt.Println("参数错误：文本必须用双引号包裹（例如 \"XYZ\"）")
		return
	}
	// 提取引号内的文本（支持包含换行符 \n）
	content := textArg[1 : len(textArg)-1]

	// 5. 执行插入操作（调用编辑器的 Insert 方法）
	activeEditor.Insert(line, col, content)
	fmt.Printf("已在 %d:%d 位置插入文本：%s\n", line, col, content)
}

func _delete(ws *workspace.Workspace, parts []string) {
	// 1. 校验活动文件是否存在
	activeEditor := ws.GetActiveEditor()
	if activeEditor == nil {
		fmt.Println("错误：没有打开的文件，请先使用 load 命令加载文件")
		return
	}

	// 2. 校验参数数量（必须包含 <line:col> 和 <len> 两个参数）
	if len(parts) != 3 {
		fmt.Println("参数错误：格式为 delete <line:col> <len>（例如 delete 1:7 5）")
		return
	}

	// 3. 解析位置参数 <line:col>
	posStr := parts[1]
	var line, col int
	posParts := strings.Split(posStr, ":")
	if len(posParts) != 2 {
		fmt.Println("参数错误：位置格式应为 line:col（例如 1:7）")
		return
	}
	// 行号必须为正整数
	line, err := strconv.Atoi(posParts[0])
	if err != nil || line < 1 {
		fmt.Println("参数错误：行号必须为正整数")
		return
	}
	// 列号必须为正整数
	col, err = strconv.Atoi(posParts[1])
	if err != nil || col < 1 {
		fmt.Println("参数错误：列号必须为正整数")
		return
	}

	// 4. 解析删除长度 <len>
	lenStr := parts[2]
	length, err := strconv.Atoi(lenStr)
	if err != nil || length < 1 {
		fmt.Println("参数错误：删除长度必须为正整数")
		return
	}

	// 5. 执行删除操作（调用编辑器的 Delete 方法）
	// 编辑器内部会处理：行号/列号越界、删除长度超出行尾等异常
	activeEditor.Delete(line, col, length)
	fmt.Printf("已从 %d:%d 位置删除 %d 个字符\n", line, col, length)
}

func _replace(ws *workspace.Workspace, parts []string) {
	// 1. 校验活动文件是否存在
	activeEditor := ws.GetActiveEditor()
	if activeEditor == nil {
		fmt.Println("错误：没有打开的文件，请先使用 load 命令加载文件")
		return
	}

	// 2. 校验参数数量（必须包含 <line:col>、<len>、"text" 三个参数）
	if len(parts) < 4 {
		fmt.Println("参数错误：格式为 replace <line:col> <len> \"text\"（例如 replace 1:1 4 \"slow\"）")
		return
	}

	// 3. 解析位置参数 <line:col>
	posStr := parts[1]
	var line, col int
	//posParts按 ":" 分割行号和列号
	posParts := strings.Split(posStr, ":")
	if len(posParts) != 2 {
		fmt.Println("参数错误：位置格式应为 line:col（例如 1:1）")
		return
	}
	// 行号必须为正整数
	line, err := strconv.Atoi(posParts[0])
	if err != nil || line < 1 {
		fmt.Println("参数错误：行号必须为正整数")
		return
	}
	// 列号必须为正整数
	col, err = strconv.Atoi(posParts[1])
	if err != nil || col < 1 {
		fmt.Println("参数错误：列号必须为正整数")
		return
	}

	// 4. 解析删除删除长度 <len>
	lenStr := parts[2]
	length, err := strconv.Atoi(lenStr)
	if err != nil || length < 1 {
		fmt.Println("参数错误：删除长度必须为正整数")
		return
	}

	// 5. 解析替换文本（支持带空格和空字符串）
	// 文本参数从 parts[3] 开始，合并所有后续片段
	textParts := parts[3:]
	textArg := strings.Join(textParts, " ")
	// 校验文本是否用双引号包裹（空字符串需表示为 ""）
	if len(textArg) < 2 || textArg[0] != '"' || textArg[len(textArg)-1] != '"' {
		fmt.Println("参数错误：替换文本必须用双引号包裹（例如 \"slow\" 或 \"\"）")
		return
	}
	// 提取引号内的文本（支持空字符串）
	content := textArg[1 : len(textArg)-1]

	// 6. 执行替换操作（调用编辑器的 Replace 方法）
	// 编辑器内部会先执行 delete 再执行 insert，处理各类异常
	activeEditor.Replace(line, col, length, content)
	fmt.Printf("已从 %d:%d 位置替换 %d 个字符为：%s\n", line, col, length, content)
}
