package treeview

import (
	"strings"
)

// TreeView renders any TreeNode structure to a text tree format
// Inspired by VS Code's TreeView but simplified for console output
type TreeView struct {
	root TreeNode
}

// NewTreeView creates a new TreeView with the given root node
func NewTreeView(root TreeNode) *TreeView {
	return &TreeView{root: root}
}

// Render generates the tree structure as a formatted string
func (tv *TreeView) Render() string {
	if tv.root == nil {
		return ""
	}
	
	var builder strings.Builder
	tv.renderNode(tv.root, "", true, &builder)
	return builder.String()
}

// renderNode recursively renders a node and its children
func (tv *TreeView) renderNode(node TreeNode, prefix string, isLast bool, builder *strings.Builder) {
	// Write the node label with appropriate prefix
	builder.WriteString(prefix)
	
	if prefix != "" {
		if isLast {
			builder.WriteString("└── ")
		} else {
			builder.WriteString("├── ")
		}
	}
	
	builder.WriteString(node.GetLabel())
	builder.WriteString("\n")
	
	// Process children if this is not a leaf node
	if !node.IsLeaf() {
		children := node.GetChildren()
		for i, child := range children {
			isChildLast := i == len(children)-1
			
			// Calculate the prefix for child nodes
			var childPrefix string
			if prefix == "" {
				// For root node, use simple spacing
				childPrefix = "    "
			} else {
				if isLast {
					childPrefix = prefix + "    "
				} else {
					childPrefix = prefix + "│   "
				}
			}
			
			tv.renderNode(child, childPrefix, isChildLast, builder)
		}
	}
}
