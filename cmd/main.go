package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/gen"
	"github.com/codecrafters-io/interpreter-starter-go/tokens"
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

// Parser start===>
type Parser struct {
	tokens  []*tokens.Token
	current int
}

func NewParser(tokens_ []*tokens.Token) *Parser {
	return &Parser{
		tokens:  tokens_,
		current: 0,
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if (command != "tokenize") && (command != "gen-ast") {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	if command == "gen-ast" {
		gen.GenerateAST()
	} else {
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		tokens.CreateKeyWords()

		source := string(fileContents)

		scanner := NewScanner(source)
		tokens := scanner.ScanTokens()

		for _, token := range tokens {
			fmt.Println(token.String())
		}
	}

	// Exit with code 65 if any errors occurred
	if hadError {
		os.Exit(65)
	}
}
