package tokens

import "fmt"

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

var Keywords map[string]TokenType

func CreateKeyWords() {
	Keywords = make(map[string]TokenType)

	Keywords["and"] = AND
	Keywords["class"] = CLASS
	Keywords["else"] = ELSE
	Keywords["false"] = FALSE
	Keywords["fun"] = FUN
	Keywords["for"] = FOR
	Keywords["if"] = IF
	Keywords["nil"] = NIL
	Keywords["or"] = OR
	Keywords["print"] = PRINT
	Keywords["return"] = RETURN
	Keywords["super"] = SUPER
	Keywords["this"] = THIS
	Keywords["true"] = TRUE
	Keywords["var"] = VAR
	Keywords["while"] = WHILE
}
