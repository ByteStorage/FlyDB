package engine

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
)

// ASTNode is the abstract syntax tree node.
type ASTNode interface {
	// String returns the string representation of the node.
	String() string
	Only() bool
}

type SQLParser interface {
	Parse(query string) (ASTNode, error)
}

type SelectStatement struct {
	Columns []string
	From    string
	Where   Condition
}

type InsertStatement struct {
	Columns []string
	From    string
}

func (s SelectStatement) String() string {
	return fmt.Sprintf("select %v from %s where %s", s.Columns, s.From, s.Where)
}

func (s SelectStatement) Only() bool {
	return true
}

func NewSQLParser(sql string) SQLParser {
	return &parser{
		sql: sql,
	}
}

type Condition struct {
	Column   string
	Operator string
	Value    interface{}
}

func (c *Condition) String() string {
	return fmt.Sprintf("%s %s %v", c.Column, c.Operator, c.Value)
}

type parser struct {
	sql string
}

func (p *parser) Parse(query string) (ASTNode, error) {
	statement, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}
	switch statement := statement.(type) {
	case *sqlparser.Select:
		return p.parseSelect(statement)
	case *sqlparser.Insert:
		return p.parseInsert(statement)
	default:
		return nil, fmt.Errorf("unsupported statement type: %T", statement)
	}
}

func (p *parser) parseSelect(statement *sqlparser.Select) (*SelectStatement, error) {
	var columns []string
	for _, column := range statement.SelectExprs {
		switch column := column.(type) {
		case *sqlparser.AliasedExpr:
			columns = append(columns, column.Expr.(*sqlparser.ColName).Name.String())
		default:
			return nil, fmt.Errorf("unsupported select expression type: %T", column)
		}
	}
	table := statement.From[0].(*sqlparser.AliasedTableExpr).Expr.(*sqlparser.TableName).Name.String()
	var condition Condition
	if statement.Where != nil {
		condition = p.parseCondition(statement.Where.Expr)
	}
	return &SelectStatement{
		Columns: columns,
		From:    table,
		Where:   condition,
	}, nil
}

func (p *parser) parseInsert(statement *sqlparser.Insert) (*SelectStatement, error) {
	var columns []string
	for _, column := range statement.Columns {
		columns = append(columns, column.String())
	}
	table := statement.Table.Name.String()
	return &SelectStatement{
		Columns: columns,
		From:    table,
	}, nil
}

func (p *parser) parseCondition(expr sqlparser.Expr) Condition {
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return Condition{
			Column:   expr.Left.(*sqlparser.ColName).Name.String(),
			Operator: expr.Operator,
			Value:    expr.Right,
		}
	default:
		panic(fmt.Sprintf("unsupported expression type: %T", expr))
	}
}
