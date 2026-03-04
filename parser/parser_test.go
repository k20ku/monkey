package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foo = 33;
	`
	// lexer
	l := lexer.New(input)
	fmt.Printf("l Lexer = %#v\n", l)
	// parser
	p := New(l)
	fmt.Printf("p Parser = %#v\n", p)
	fmt.Printf("l Lexer = %#v\n\n", l)

	// parse!
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// must be {x-let, y-let, foobar-let}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}

}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Erros()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	// stop the test exceution (Fail Fast!)
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	/*
		type Node = Statement | Expression
		type Program = Statement[]
		type Statement = LetStatement | ...
	*/
	// **s Statement is truely LetStatement?**
	// 1. Token
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	// 2. Can s Statement be casted to s LetStatement?
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	// Name
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not  '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not %s, got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}
	// Value
	// TODO: Value Validation

	return true
}
