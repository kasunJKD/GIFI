package main

import (
	"fmt"
	"os"
	"strconv"
)

var hadError bool

type TokenType string

const (
	//single character tokens
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = "COMMA"
	DOT         TokenType = "DOT"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	SEMICOLON   TokenType = "SEMICOLON"
	SLASH       TokenType = "SLASH"
	STAR        TokenType = "STAR"

	// One or two character tokens.
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"

	// Literals.
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	// Keywords.
	AND    TokenType = "AND"
	CLASS  TokenType = "CLASS"
	ELSE   TokenType = "ELSE"
	FALSE  TokenType = "FALSE"
	FUN    TokenType = "FUN"
	FOR    TokenType = "FOR"
	IF     TokenType = "IF"
	NIL    TokenType = "NIL"
	OR     TokenType = "OR"
	PRINT  TokenType = "PRINT"
	RETURN TokenType = "RETURN"
	SUPER  TokenType = "SUPER"
	THIS   TokenType = "THIS"
	TRUE   TokenType = "TRUE"
	VAR    TokenType = "VAR"
	WHILE  TokenType = "WHILE"
	EOF    TokenType = "EOF"
)

var keywords map[string]TokenType

func createKeyWords() {
	keywords = make(map[string]TokenType)

	keywords["and"] = AND
	keywords["class"] = CLASS
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["fun"] = FUN
	keywords["for"] = FOR
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["true"] = TRUE
	keywords["var"] = VAR
	keywords["while"] = WHILE
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	var literal string
	if t.Literal == nil {
		literal = "null"
	} else {
		literal = fmt.Sprintf("%v", t.Literal)
	}
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, literal)
}

// scanner->start
type Scanner struct {
	Tokens []*Token
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
func (s *Scanner) ScanTokens() []*Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, &Token{
		Type:   EOF,
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
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '=':
		var enumval TokenType
		if s.match('=') {
			enumval = EQUAL_EQUAL
		} else {
			enumval = EQUAL
		}
		s.addToken(enumval, nil)
	case '!':
		var enumval TokenType
		if s.match('=') {
			enumval = BANG_EQUAL
		} else {
			enumval = BANG
		}
		s.addToken(enumval, nil)
	case '<':
		var enumval TokenType
		if s.match('=') {
			enumval = LESS_EQUAL
		} else {
			enumval = LESS
		}
		s.addToken(enumval, nil)
	case '>':
		var enumval TokenType
		if s.match('=') {
			enumval = GREATER_EQUAL
		} else {
			enumval = GREATER
		}
		s.addToken(enumval, nil)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
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
	s.addToken(STRING, value)
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
	s.addToken(NUMBER, formattedValue)
}

func (s *Scanner) identifer() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.Source[s.start:s.current]
	tokenType := keywords[text]
	if tokenType == "" {
		tokenType = IDENTIFIER
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

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, &Token{
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

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	createKeyWords()

	source := string(fileContents)

	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token.String())
	}

	// Exit with code 65 if any errors occurred
	if hadError {
		os.Exit(65)
	}
}
