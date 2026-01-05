#!/bin/bash

# 测试脚本：演示树形结构显示功能

echo "==================================="
echo "树形结构显示功能测试"
echo "==================================="
echo ""

# 编译项目
echo "1. 编译项目..."
go build -o lab1 || exit 1
echo "   ✓ 编译成功"
echo ""

# 测试目录树
echo "2. 测试目录树显示（dir-tree 命令）"
echo "   命令：dir-tree files"
echo "   输出："
echo "-----------------------------------"
echo "dir-tree files" | ./lab1 | grep -A 20 ">" | tail -n +2 | head -n 10
echo "-----------------------------------"
echo ""

# 测试 XML 树
echo "3. 测试 XML 树显示（xml-tree 命令）"
echo "   命令：xml-tree project.xml"
echo "   输出："
echo "-----------------------------------"
echo "xml-tree project.xml" | ./lab1 | grep -A 50 ">" | tail -n +2 | head -n 30
echo "-----------------------------------"
echo ""

# 测试根目录树（简化版）
echo "4. 测试当前目录树显示"
echo "   命令：dir-tree treeview"
echo "   输出："
echo "-----------------------------------"
echo "dir-tree treeview" | ./lab1 | grep -A 20 ">" | tail -n +2 | head -n 10
echo "-----------------------------------"
echo ""

echo "==================================="
echo "测试完成！"
echo "==================================="
echo ""
echo "说明："
echo "- dir-tree 命令：显示文件系统目录结构"
echo "- xml-tree 命令：显示 XML 文档树结构"
echo "- 两个命令都使用了适配器模式，共享相同的 TreeView 渲染逻辑"
echo ""
