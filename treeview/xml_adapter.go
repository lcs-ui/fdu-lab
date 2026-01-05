package treeview

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// XMLNode adapts an XML element to the TreeNode interface
type XMLNode struct {
	name       string
	attributes []xml.Attr
	children   []TreeNode
	textValue  string
	isLeaf     bool
}

// NewXMLNodeFromFile creates an XML tree from a file
func NewXMLNodeFromFile(path string) (*XMLNode, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return NewXMLNodeFromReader(file)
}

// NewXMLNodeFromReader creates an XML tree from an io.Reader
func NewXMLNodeFromReader(reader io.Reader) (*XMLNode, error) {
	decoder := xml.NewDecoder(reader)
	
	// Parse the XML document
	var root *XMLNode
	var stack []*XMLNode
	
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		
		switch t := token.(type) {
		case xml.StartElement:
			node := &XMLNode{
				name:       t.Name.Local,
				attributes: t.Attr,
				children:   make([]TreeNode, 0),
				isLeaf:     false,
			}
			
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, node)
			} else {
				root = node
			}
			
			stack = append(stack, node)
			
		case xml.EndElement:
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				if len(current.children) == 0 && current.textValue == "" {
					current.isLeaf = true
				}
				stack = stack[:len(stack)-1]
			}
			
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && len(stack) > 0 {
				current := stack[len(stack)-1]
				current.textValue = text
			}
		}
	}
	
	return root, nil
}

// GetLabel implements TreeNode interface
func (x *XMLNode) GetLabel() string {
	label := x.name
	
	// Add attributes if present
	if len(x.attributes) > 0 {
		attrStrs := make([]string, 0, len(x.attributes))
		for _, attr := range x.attributes {
			attrStrs = append(attrStrs, fmt.Sprintf("%s=\"%s\"", attr.Name.Local, attr.Value))
		}
		label = fmt.Sprintf("%s [%s]", label, strings.Join(attrStrs, ", "))
	}
	
	// Add text content if present
	if x.textValue != "" {
		label = fmt.Sprintf("%s: %s", label, x.textValue)
	}
	
	return label
}

// GetChildren implements TreeNode interface
func (x *XMLNode) GetChildren() []TreeNode {
	return x.children
}

// IsLeaf implements TreeNode interface
func (x *XMLNode) IsLeaf() bool {
	return x.isLeaf
}

// XMLTreeAdapter adapts XML documents to TreeView
// This is the Adapter that connects XML structures to our TreeView
type XMLTreeAdapter struct {
	filePath string
}

// NewXMLTreeAdapter creates a new adapter for the given XML file path
func NewXMLTreeAdapter(filePath string) *XMLTreeAdapter {
	return &XMLTreeAdapter{filePath: filePath}
}

// GetRootNode returns the root TreeNode for this XML document
func (xta *XMLTreeAdapter) GetRootNode() (TreeNode, error) {
	return NewXMLNodeFromFile(xta.filePath)
}
