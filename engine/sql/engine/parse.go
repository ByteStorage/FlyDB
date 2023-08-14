package engine

import "fmt"

// ASTNode is the abstract syntax tree node.
type ASTNode interface {
	// NodeType returns node type.
	NodeType() string
	// Children returns child nodes.
	Children() []ASTNode
	// String returns the string representation of the node.
	String() string
}

type SQLParser interface {
	Parse(query string) (ASTNode, error)
}

type SelectStatement struct {
	Columns []string
	From    string
	Where   Condition
}

func (s *SelectStatement) String() string {
	panic("implement me")
}

type Condition struct {
	Column   string
	Operator string
	Value    interface{}
}

func (c *Condition) String() string {
	return fmt.Sprintf("%s %s %v", c.Column, c.Operator, c.Value)
}

type SimpleSQLParser struct{}

func (p *SimpleSQLParser) Parse(query string) (ASTNode, error) {
	panic("implement me")
}
