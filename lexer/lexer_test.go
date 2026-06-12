package lexer

import (
	"testing"

	"github.com/k20ku/monkey/token"
)

func TestOnlySymbols(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf(
				"test[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type,
			)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"test[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal,
			)
		}
	}
}

type expectedToken struct {
	Type    token.TokenType
	Literal string
}

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []expectedToken
	}{
		{
			"let five = 5;",
			[]expectedToken{
				{token.LET, "let"},
				{token.IDENT, "five"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"let ten = 10;",
			[]expectedToken{
				{token.LET, "let"},
				{token.IDENT, "ten"},
				{token.ASSIGN, "="},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"let add = fn(x, y) {x + y;};",
			[]expectedToken{
				{token.LET, "let"},
				{token.IDENT, "add"},
				{token.ASSIGN, "="},
				{token.FUNCTION, "fn"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.COMMA, ","},
				{token.IDENT, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.IDENT, "x"},
				{token.PLUS, "+"},
				{token.IDENT, "y"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"let result = add(five, ten);",
			[]expectedToken{
				{token.LET, "let"},
				{token.IDENT, "result"},
				{token.ASSIGN, "="},
				{token.IDENT, "add"},
				{token.LPAREN, "("},
				{token.IDENT, "five"},
				{token.COMMA, ","},
				{token.IDENT, "ten"},
				{token.RPAREN, ")"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"let area55 = 55;",
			[]expectedToken{
				{token.LET, "let"},
				{token.IDENT, "area55"},
				{token.ASSIGN, "="},
				{token.INT, "55"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"!-/*5;",
			[]expectedToken{
				{token.BANG, "!"},
				{token.MINUS, "-"},
				{token.SLASH, "/"},
				{token.ASTERISK, "*"},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"5<10>5",
			[]expectedToken{
				{token.INT, "5"},
				{token.LT, "<"},
				{token.INT, "10"},
				{token.GT, ">"},
				{token.INT, "5"},
				{token.EOF, ""},
			},
		},
		{
			"if (5 > 10) { return true; } else { return false; }",
			[]expectedToken{
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.INT, "5"},
				{token.GT, ">"},
				{token.INT, "10"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.TRUE, "true"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.ELSE, "else"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.FALSE, "false"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			"10 == 10;",
			[]expectedToken{
				{token.INT, "10"},
				{token.EQ, "=="},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"10 != 9;",
			[]expectedToken{
				{token.INT, "10"},
				{token.NOT_EQ, "!="},
				{token.INT, "9"},
				{token.SEMICOLON, ";"},
			},
		},
		{
			`"foobar"`,
			[]expectedToken{
				{token.STRING, "foobar"},
				{token.EOF, ""},
			},
		},
		{
			`"foo bar"`,
			[]expectedToken{
				{token.STRING, "foo bar"},
			},
		},
		{
			`"Hello \t \n \\ \"World\""`,
			[]expectedToken{
				{token.STRING, "Hello \t \n \\ \"World\""},
			},
		},
		{
			`puts("Hello");sleep(2);puts("\rWorld!")`,
			[]expectedToken{
				{token.IDENT, "puts"},
				{token.LPAREN, "("},
				{token.STRING, "Hello"},
				{token.RPAREN, ")"},
				{token.SEMICOLON, ";"},
				{token.IDENT, "sleep"},
				{token.LPAREN, "("},
				{token.INT, "2"},
				{token.RPAREN, ")"},
				{token.SEMICOLON, ";"},
				{token.IDENT, "puts"},
				{token.LPAREN, "("},
				{token.STRING, "\rWorld!"},
				{token.RPAREN, ")"},
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input)

			for _, expectedTok := range tt.expected {
				tok := l.NextToken()

				if tok.Type != expectedTok.Type {
					t.Fatalf(
						"tokentype wrong. expected=%q, got=%q",
						expectedTok.Type, tok.Type,
					)
				}

				if tok.Literal != expectedTok.Literal {
					t.Fatalf(
						"tokenliteral wrong. expected=%q, got=%q",
						expectedTok.Literal, tok.Literal,
					)
				}
			}
		})
	}
}
