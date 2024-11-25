package  gen
import github.com/codecrafters-io/interpreter-starter-go/tokens
type Expr interface {
	Accept(visitor VisitorExpr) interface{}
}
type VisitorExpr interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
}
type Binary struct  {
	 Left Expr
	 Right Expr
	 Operator Token
}
func NewBinary(Left Expr, Right Expr, Operator Token) *Binary{
	return & Binary {
		 Left :	 Left ,
		 Right :	 Right ,
		 Operator :	 Operator ,
	}
}
func (a *Binary) Accept(v VisitorExpr) interface{} {
	return v.VisitBinaryExpr(a)
}
type Unary struct  {
	 Operator Token
	 Right Expr
}
func NewUnary(Operator Token, Right Expr) *Unary{
	return & Unary {
		 Operator :	 Operator ,
		 Right :	 Right ,
	}
}
func (a *Unary) Accept(v VisitorExpr) interface{} {
	return v.VisitUnaryExpr(a)
}
type Grouping struct  {
	 Expression Expr
}
func NewGrouping(Expression Expr) *Grouping{
	return & Grouping {
		 Expression :	 Expression ,
	}
}
func (a *Grouping) Accept(v VisitorExpr) interface{} {
	return v.VisitGroupingExpr(a)
}
type Literal struct  {
	 Value interface{}
}
func NewLiteral(Value interface{}) *Literal{
	return & Literal {
		 Value :	 Value ,
	}
}
func (a *Literal) Accept(v VisitorExpr) interface{} {
	return v.VisitLiteralExpr(a)
}
