package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go-intepreter/gen"
	"go-intepreter/tokens"
)

var hadError bool

// scanner->start
type Scanner struct {
	Tokens []*tokens.Token
	Source string

	start   int
	current int
	line    int
}

// NewScanner creates a new Scanner instance
func NewScanner(source string) *Scanner {
	return &Scanner{
		Source:  source,
		line:    1,
		start:   0,
		current: 0,
	}
}

// ScanTokens scans all tokens in the source
func (s *Scanner) ScanTokens() []*tokens.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, &tokens.Token{
		Type:   tokens.EOF,
		Lexeme: "",
		Line:   s.line,
	})

	return s.Tokens
}

// isAtEnd checks if the scanner has reached the end of the source
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

// scanToken scans a single token
func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(tokens.LEFT_PAREN, nil)
	case ')':
		s.addToken(tokens.RIGHT_PAREN, nil)
	case '{':
		s.addToken(tokens.LEFT_BRACE, nil)
	case '}':
		s.addToken(tokens.RIGHT_BRACE, nil)
	case ',':
		s.addToken(tokens.COMMA, nil)
	case '.':
		s.addToken(tokens.DOT, nil)
	case '-':
		s.addToken(tokens.MINUS, nil)
	case '+':
		s.addToken(tokens.PLUS, nil)
	case ';':
		s.addToken(tokens.SEMICOLON, nil)
	case '*':
		s.addToken(tokens.STAR, nil)
	case '=':
		var enumval tokens.TokenType
		if s.match('=') {
			enumval = tokens.EQUAL_EQUAL
		} else {
			enumval = tokens.EQUAL
		}
		s.addToken(enumval, nil)
	case '!':
		var enumval tokens.TokenType
		if s.match('=') {
			enumval = tokens.BANG_EQUAL
		} else {
			enumval = tokens.BANG
		}
		s.addToken(enumval, nil)
	case '<':
		var enumval tokens.TokenType
		if s.match('=') {
			enumval = tokens.LESS_EQUAL
		} else {
			enumval = tokens.LESS
		}
		s.addToken(enumval, nil)
	case '>':
		var enumval tokens.TokenType
		if s.match('=') {
			enumval = tokens.GREATER_EQUAL
		} else {
			enumval = tokens.GREATER
		}
		s.addToken(enumval, nil)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tokens.SLASH, nil)
		}
	case '\n':
		s.line++
		break
	case '\t':
		break
	case ' ':
		break
	case '\r':
		break
	case '"':
		s.string()
		break
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifer()
		} else {
			error(s.line, "Unexpected character:", string(c))
			os.Exit(1) // Terminate on unexpected characters
		}
	}
}

func (s *Scanner) advance() byte {
	c := s.Source[s.current]
	s.current++
	return c
}

func (s *Scanner) match(nextExpected byte) bool {
	if s.isAtEnd() {
		return false
	}
	c := s.Source[s.current]
	if c != nextExpected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\n'
	} else {
		return s.Source[s.current]
	}
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		error(s.line, "Unterminated string.", "")
		return
	}

	s.advance()

	value := s.Source[s.start+1 : s.current-1]
	s.addToken(tokens.STRING, value)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.nextPeek()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	convvalue, _ := strconv.ParseFloat(s.Source[s.start:s.current], 64)
	// Format the value
	var formattedValue string
	if convvalue == float64(int64(convvalue)) {
		// If the value is an integer, format with one decimal place
		formattedValue = fmt.Sprintf("%.1f", convvalue)
	} else {
		// Otherwise, use the full precision as is
		formattedValue = fmt.Sprintf("%v", convvalue)
	}

	// Add the formatted value as a token
	s.addToken(tokens.NUMBER, formattedValue)
}

func (s *Scanner) identifer() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.Source[s.start:s.current]
	tokenType := tokens.Keywords[text]
	if tokenType == "" {
		tokenType = tokens.IDENTIFIER
	}

	s.addToken(tokenType, nil)
}

func (s *Scanner) isAlpha(c byte) bool {
	return ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_'))
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return (s.isAlpha(c) || s.isDigit(c))
}

func (s *Scanner) nextPeek() byte {
	if s.current+1 <= len(s.Source) {
		return s.Source[s.current+1]
	}
	return '"'
}

func (s *Scanner) addToken(tokenType tokens.TokenType, literal interface{}) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, &tokens.Token{
		Type:    tokenType,
		Lexeme:  text,
		Literal: literal,
		Line:    s.line,
	})
}

//scanner -> end

// error logs->start
func error(line int, message string, value string) {
	report(line, "", message, value)
	hadError = true
}

// report formats and logs the error message to stderr.
func report(line int, where string, message string, value string) {
	if value != "" {
		fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s %s\n", line, where, message, value)

	} else {
		fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	}
}

// AST expr
// Expr section
// type Expr interface {
// 	accpet(visiter VisitorExpr) interface{}
// }
//
// type VisitorExpr interface {
// 	VisitBinaryExpr(expr *Binary) interface{}
// }

// visitor type definitons for interface @use this to access !!!!!
// type VisitorType struct{}
//
// func (v VisitorType) VisitBinaryExpr(expr *Binary) interface{} {
//
// }

// binary expr section
// type Binary struct {
// 	Left     Expr
// 	Right    Expr
// 	Operator tokens.Token
// }
//
// func NewBinary(left Expr, right Expr, op tokens.Token) *Binary {
// 	return &Binary{
// 		Left:     left,
// 		Right:    right,
// 		Operator: op,
// 	}
// }

// binary accpet interface base
// func (b *Binary) accpet(v VisitorExpr) interface{} {
// 	return v.VisitBinaryExpr(b)
// }

// visitor type definitons for interface @use this to access !!!!!
type ASTPrinter struct{}

func (a *ASTPrinter) VisitBinaryExpr(expr *gen.Binary) string {
	return a.Paranthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *ASTPrinter) VisitGroupingExpr(expr *gen.Grouping) string {
	return a.Paranthesize("group", expr.Expression)
}

func (a *ASTPrinter) VisitLiteralExpr(expr *gen.Literal) string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *ASTPrinter) VisitUnaryExpr(expr *gen.Unary) string {
	return a.Paranthesize(expr.Operator.Lexeme, expr.Right)
}

func (v *ASTPrinter) Paranthesize(name string, exprs ...gen.Expr) string {
    var builder strings.Builder

    builder.WriteString("(" + name)
    for _, expr := range exprs {
        if expr != nil {
            builder.WriteString(" ")
            builder.WriteString(expr.Accept(v))
        } else {
            builder.WriteString(" <nil> ") // Handle nil expressions gracefully
        }
    }
    builder.WriteString(")")

    return builder.String()
}
// Parser start===>
type Parser struct {
	tokens  []*tokens.Token
	current int
}

func NewParser(tokens []*tokens.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) expression() gen.Expr {
	return p.equality()
}

func (p *Parser) equality() gen.Expr {
	expr := p.comparison()

	for p.match(tokens.BANG_EQUAL, tokens.EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparison()
		expr = gen.NewBinary(expr, right, op)
	}
	return expr
}

func (p *Parser) comparison() gen.Expr {
	expr := p.term()

	for p.match(tokens.GREATER, tokens.GREATER_EQUAL, tokens.LESS, tokens.LESS_EQUAL) {
		op := p.previous()
		right := p.term()
		expr = gen.NewBinary(expr, right, op)
	}
	return expr
}

func (p *Parser) term() gen.Expr {
	expr := p.factor()

	for p.match(tokens.MINUS, tokens.PLUS) {
		op := p.previous()
		right := p.factor()
		expr = gen.NewBinary(expr, right, op)
	}
	return expr
}

func (p *Parser) factor() gen.Expr {
	expr := p.unary()

	for p.match(tokens.SLASH, tokens.STAR) {
		op := p.previous()
		right := p.unary()
		expr = gen.NewBinary(expr, right, op)
	}
	return expr
}

func (p *Parser) unary() gen.Expr {
	if p.match(tokens.BANG, tokens.MINUS) {
		op := p.previous()
		right := p.unary()
		return gen.NewUnary(op, right)
	}
	return p.primary()
}

func (p *Parser) primary() gen.Expr {
	if p.match(tokens.FALSE) {
		return gen.NewLiteral(false)
	}
	if p.match(tokens.TRUE) {
		return gen.NewLiteral(true)
	}
	if p.match(tokens.NIL) {
		return gen.NewLiteral(nil)
	}
	if p.match(tokens.NUMBER, tokens.STRING) {
		return gen.NewLiteral(p.previous().Literal)
	}
	if p.match(tokens.LEFT_PAREN) {
		expr := p.expression()
		p.consume(tokens.RIGHT_PAREN, "Expect ')' after expression")
		return gen.NewGrouping(expr)
	}
	error(p.peek().Line, "Expect expression.", "")
	return nil
}

func (p *Parser) parse() gen.Expr {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Error during parsing:", r)
        }
    }()
    expr := p.expression()
    if expr == nil {
        error(p.peek().Line, "Failed to parse expression.", "")
        os.Exit(1) // Terminate parsing on error
    }
    return expr
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == tokens.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case tokens.CLASS, tokens.FUN, tokens.VAR, tokens.FOR, tokens.IF, tokens.WHILE, tokens.PRINT, tokens.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) match(types ...tokens.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() *tokens.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == tokens.EOF
}

func (p *Parser) peek() *tokens.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *tokens.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType tokens.TokenType, message string) *tokens.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	return nil
	//panic(p.error(p.peek(), message))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./your_program.sh <command> <source-file>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %s\n", err)
		os.Exit(65)
	}

	source := string(data)
	switch command {
	case "tokenize":
		scanner := NewScanner(source)
		tokens := scanner.ScanTokens()
		for _, token := range tokens {
			fmt.Println(token)
		}

case "parse":
    scanner := NewScanner(source)
    tokens := scanner.ScanTokens()
    if hadError {
        os.Exit(1) // Stop if scanning failed
    }
    parser := NewParser(tokens)
    expr := parser.parse()
    if expr != nil {
        printer := &ASTPrinter{}
        switch e := expr.(type) {
        case *gen.Binary:
            fmt.Println(printer.VisitBinaryExpr(e))
        case *gen.Grouping:
            fmt.Println(printer.VisitGroupingExpr(e))
        case *gen.Unary:
            fmt.Println(printer.VisitUnaryExpr(e))
        case *gen.Literal:
            fmt.Println(printer.VisitLiteralExpr(e))
        default:
            fmt.Println("Unknown expression type.")
        }
    } else {
        fmt.Println("Parsing failed.")
    }	

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
