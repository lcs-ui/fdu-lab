package treeview

// TreeNode represents a node in a tree structure
// This is the core interface that all adapters must implement
type TreeNode interface {
	// GetLabel returns the display text for this node
	GetLabel() string
	
	// GetChildren returns child nodes (nil if this is a leaf node)
	GetChildren() []TreeNode
	
	// IsLeaf returns true if this node has no children
	IsLeaf() bool
}
