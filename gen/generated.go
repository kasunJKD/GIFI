package gen

import tokens "go-intepreter/tokens"

type Expr interface {
	Accept(visitor VisitorExpr) string
}
type VisitorExpr interface {
	VisitBinaryExpr(expr *Binary) string
	VisitUnaryExpr(expr *Unary) string
	VisitGroupingExpr(expr *Grouping) string
	VisitLiteralExpr(expr *Literal) string
}
type Binary struct {
	Left     Expr
	Right    Expr
	Operator *tokens.Token
}

func NewBinary(Left Expr, Right Expr, Operator *tokens.Token) *Binary {
	return &Binary{
		Left:     Left,
		Right:    Right,
		Operator: Operator,
	}
}
func (a *Binary) Accept(v VisitorExpr) string {
	return v.VisitBinaryExpr(a)
}

type Unary struct {
	Operator *tokens.Token
	Right    Expr
}

func NewUnary(Operator *tokens.Token, Right Expr) *Unary {
	return &Unary{
		Operator: Operator,
		Right:    Right,
	}
}
func (a *Unary) Accept(v VisitorExpr) string {
	return v.VisitUnaryExpr(a)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) *Grouping {
	return &Grouping{
		Expression: Expression,
	}
}
func (a *Grouping) Accept(v VisitorExpr) string {
	return v.VisitGroupingExpr(a)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(Value interface{}) *Literal {
	return &Literal{
		Value: Value,
	}
}
func (a *Literal) Accept(v VisitorExpr) string {
	return v.VisitLiteralExpr(a)
}
