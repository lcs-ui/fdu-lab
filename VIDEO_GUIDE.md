# 视频录制指南 - 树形结构显示的适配器模式重构

## 视频时长：1-2分钟

## 场景一：展示代码结构（30-40秒）

### 1. 展示 treeview 模块的文件结构
```bash
# 在终端中运行
ls -la treeview/
```

应该看到：
- `tree_node.go` - TreeNode 接口定义
- `tree_view.go` - TreeView 渲染器
- `directory_adapter.go` - 目录树适配器
- `xml_adapter.go` - XML 树适配器

### 2. 快速浏览关键代码文件

**展示 tree_node.go（5秒）**
- 指出 TreeNode 接口的三个核心方法：GetLabel(), GetChildren(), IsLeaf()
- 说明：这是适配器模式的目标接口

**展示 tree_view.go（5秒）**
- 指出 TreeView 只依赖 TreeNode 接口
- 说明：渲染逻辑与具体数据源解耦

**展示 directory_adapter.go（5秒）**
- 指出 DirectoryNode 实现了 TreeNode 接口
- 说明：将文件系统目录适配为树节点

**展示 xml_adapter.go（5秒）**
- 指出 XMLNode 实现了 TreeNode 接口
- 说明：将 XML 文档结构适配为树节点

### 3. 展示 main.go 中的命令处理（10秒）
```go
case "dir-tree":
    adapter := treeview.NewDirectoryTreeAdapter(targetDir)
    rootNode, _ := adapter.GetRootNode()
    view := treeview.NewTreeView(rootNode)
    fmt.Print(view.Render())
    
case "xml-tree":
    adapter := treeview.NewXMLTreeAdapter(xmlPath)
    rootNode, _ := adapter.GetRootNode()
    view := treeview.NewTreeView(rootNode)
    fmt.Print(view.Render())
```

说明：
- 两个命令使用相同的 TreeView 渲染器
- 只需要替换不同的适配器
- 体现了适配器模式的优势

## 场景二：运行演示（40-60秒）

### 1. 编译项目（5秒）
```bash
go build -o lab1
```

### 2. 运行程序（5秒）
```bash
./lab1
```

### 3. 演示 dir-tree 命令（15秒）
```bash
> dir-tree files
```

**输出说明：**
- 展示文件目录的树形结构
- 使用 ├──、└── 等字符绘制树形连接线
- 说明：这是 DirectoryTreeAdapter 将文件系统适配为树形结构的结果

### 4. 演示 xml-tree 命令（15秒）
```bash
> xml-tree project.xml
```

**输出说明：**
- 展示 XML 文档的树形结构
- XML 元素显示为节点，属性显示在方括号中
- 文本内容显示在冒号后
- 说明：这是 XMLTreeAdapter 将 XML 结构适配为树形结构的结果

### 5. 展示更复杂的目录树（可选，10秒）
```bash
> dir-tree treeview
```

**输出说明：**
- 展示 treeview 模块自己的文件结构
- 强调：同一个渲染器，不同的数据源

### 6. 退出程序（5秒）
```bash
> exit
```

## 场景三：总结（10-20秒）

### 关键要点说明：

1. **适配器模式的实现**
   - TreeNode 作为统一接口
   - DirectoryTreeAdapter 和 XMLTreeAdapter 作为适配器
   - TreeView 作为客户端，只依赖接口

2. **设计优势**
   - 代码复用：树形渲染逻辑只写一次
   - 易于扩展：添加新数据源只需实现 TreeNode 接口
   - 解耦合：TreeView 不依赖具体数据源

3. **参考 VS Code TreeView**
   - 类似的设计思想
   - 简化为控制台输出版本

## 录制建议

### 工具选择
- 使用 OBS Studio 或类似录屏软件
- 录制终端窗口和代码编辑器

### 录制技巧
1. **清晰的窗口布局**
   - 左侧：代码编辑器（展示代码结构）
   - 右侧：终端（运行演示）

2. **语速和节奏**
   - 语速适中，重点清晰
   - 每个场景之间留 1-2 秒过渡

3. **重点突出**
   - 使用鼠标指针或高亮突出关键代码
   - 在终端输出时，可以暂停 1-2 秒让观众看清楚

### 录制脚本示例

"大家好，我来展示使用适配器模式重构树形结构显示的成果。

**[场景一]** 首先看代码结构。在 treeview 模块中，我们定义了 TreeNode 接口作为统一的树节点表示，TreeView 作为渲染器，以及两个适配器：DirectoryTreeAdapter 将文件系统适配为树形结构，XMLTreeAdapter 将 XML 文档适配为树形结构。

**[场景二]** 现在运行程序演示。首先用 dir-tree 命令显示文件目录，可以看到清晰的树形结构。然后用 xml-tree 命令显示 XML 文档，同样是树形结构，但数据来源完全不同。

**[场景三]** 这就是适配器模式的优势：一个 TreeView 渲染器，配合不同的适配器，就能显示各种数据源的树形结构。代码复用、易于扩展、高度解耦。谢谢观看！"

## 常见问题准备

Q: 为什么要用适配器模式？
A: 避免重复实现树形渲染逻辑，提高代码复用性和可维护性。

Q: 如何添加新的数据源？
A: 只需实现 TreeNode 接口，无需修改 TreeView 渲染器。

Q: 与 VS Code TreeView 的关系？
A: 参考了 VS Code 的设计理念，但简化为控制台输出版本。
