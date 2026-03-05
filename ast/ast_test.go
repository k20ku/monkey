package ast

import (
	"fmt"
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	// let myVar = anothorVar;
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anothorVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	fmt.Printf("program:\n%#v\n\n", program)

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String wrong, got=%q", program.String())
	}
}
