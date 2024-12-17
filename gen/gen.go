package gen

import (
	"fmt"
	"os"
	"strings"
)

type ExprType string

const (
	BINARY   ExprType = "Left Expr, Right Expr, Operator Token"
	GROUPING ExprType = "Expression Expr"
	LITERAL  ExprType = "Value interface{}"
	UNARY    ExprType = "Operator Token, Right Expr"
)

var exprTypeNames = map[ExprType]string{
	BINARY:   "Binary",
	UNARY:    "Unary",
	GROUPING: "Grouping",
	LITERAL:  "Literal",
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func DefineAST(outputdir string, packageName string, types []ExprType) {
	//create file
	path := outputdir
	f, err := os.Create(path)
	check(err)
	defer f.Close()

	//define base and inserts
	_, err = fmt.Fprintln(f, "package ", packageName)
	check(err)

	line := fmt.Sprintln(`import "github.com/codecrafters-io/interpreter-starter-go/tokens"`)
	//imports
	fmt.Fprintln(f, line)
	//expr interface
	_, err = fmt.Fprintln(f, "type Expr interface {")
	check(err)
	_, err = fmt.Fprintln(f, "	Accept(visitor VisitorExpr) interface{}")
	check(err)
	_, err = fmt.Fprintln(f, "}")
	check(err)

	//VisitorExpr
	_, err = fmt.Fprintln(f, "type VisitorExpr interface {")
	check(err)
	//body
	err = VisitorExprBody(f, types)
	check(err)

	_, err = fmt.Fprintln(f, "}")
	check(err)

	//typedefs and constructors
	err = TypeConstructs(f, types)
	check(err)

}

func VisitorExprBody(f *os.File, types []ExprType) error {
	for _, val := range types {
		// Construct the line using fmt.Sprintf
		line := fmt.Sprintf(
			"	Visit%sExpr(expr *%s) interface{}",
			exprTypeNames[val],
			strings.TrimSpace(exprTypeNames[val]),
		)

		// Write the line to the file
		_, err := fmt.Fprintln(f, line)
		if err != nil {
			return err
		}
	}

	return nil
}

func TypeConstructs(f *os.File, types []ExprType) error {
	for _, val := range types {
		_, err := fmt.Fprintln(f, "type", strings.TrimSpace(exprTypeNames[val]), "struct", " {")
		if err != nil {
			return err
		}

		splitarr := strings.Split(string(val), ",")
		for _, s := range splitarr {
			_, err = fmt.Fprintln(f, "	", strings.TrimSpace(s))
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintln(f, "}")
		if err != nil {
			return err
		}
		line := fmt.Sprintf("func New%s(%s) *%s{", exprTypeNames[val], strings.TrimSpace(string(val)), exprTypeNames[val])
		_, err = fmt.Fprintln(f, line)
		_, err = fmt.Fprintln(f, "	return &", exprTypeNames[val], "{")
		if err != nil {
			return err
		}
		for _, r := range splitarr {
			sp := strings.Split(strings.TrimSpace(r), " ")
			_, err = fmt.Fprintln(f, "		", sp[0], ":	", sp[0], ",")
		}
		_, err = fmt.Fprintln(f, "	}")
		_, err = fmt.Fprintln(f, "}")
		if err != nil {
			return err
		}

		line = fmt.Sprintf("func (a *%s) Accept(v VisitorExpr) interface{} {", strings.TrimSpace(exprTypeNames[val]))
		_, err = fmt.Fprintln(f, line)
		if err != nil {
			return err
		}
		line = fmt.Sprintf(
			"	return v.Visit%sExpr(a)",
			strings.TrimSpace(exprTypeNames[val]),
		)
		_, err = fmt.Fprintln(f, line)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(f, "}")
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateAST() {
	println("Generating AST")
	types := []ExprType{BINARY, UNARY, GROUPING, LITERAL}
	DefineAST("./gen/generated.go", "gen", types)
}
