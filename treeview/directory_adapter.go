package treeview

import (
	"os"
	"path/filepath"
)

// DirectoryNode adapts a file system directory to the TreeNode interface
type DirectoryNode struct {
	path     string
	name     string
	children []TreeNode
	isLeaf   bool
}

// NewDirectoryNode creates a new DirectoryNode from a file system path
func NewDirectoryNode(path string) (*DirectoryNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	
	node := &DirectoryNode{
		path:   path,
		name:   filepath.Base(path),
		isLeaf: !info.IsDir(),
	}
	
	// If it's a directory, load its children
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		
		node.children = make([]TreeNode, 0, len(entries))
		for _, entry := range entries {
			childPath := filepath.Join(path, entry.Name())
			childNode, err := NewDirectoryNode(childPath)
			if err != nil {
				continue // Skip entries that can't be read
			}
			node.children = append(node.children, childNode)
		}
	}
	
	return node, nil
}

// GetLabel implements TreeNode interface
func (d *DirectoryNode) GetLabel() string {
	return d.name
}

// GetChildren implements TreeNode interface
func (d *DirectoryNode) GetChildren() []TreeNode {
	return d.children
}

// IsLeaf implements TreeNode interface
func (d *DirectoryNode) IsLeaf() bool {
	return d.isLeaf
}

// DirectoryTreeAdapter adapts file system directories to TreeView
// This is the Adapter that connects the file system to our TreeView
type DirectoryTreeAdapter struct {
	rootPath string
}

// NewDirectoryTreeAdapter creates a new adapter for the given directory path
func NewDirectoryTreeAdapter(rootPath string) *DirectoryTreeAdapter {
	return &DirectoryTreeAdapter{rootPath: rootPath}
}

// GetRootNode returns the root TreeNode for this directory
func (dta *DirectoryTreeAdapter) GetRootNode() (TreeNode, error) {
	return NewDirectoryNode(dta.rootPath)
}
