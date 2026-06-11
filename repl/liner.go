package repl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/k20ku/monkey/lexer"
	"github.com/k20ku/monkey/token"
	liner "github.com/peterh/liner"
)

const (
	promptDefault  = "> "
	promptContinue = ". "
	indent         = "  "
)

type contLiner struct {
	*liner.State
	buffer string
	depth  int
	prompt string
}

func newContLiner() *contLiner {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	return &contLiner{State: line}
}

func (cl *contLiner) promptStringWithDepth(str string, depth int) string {
	var pr bytes.Buffer

	if str != "" {
		pr.WriteByte('(')
		pr.WriteString(str)
		pr.WriteByte(')')
	}

	if cl.buffer != "" {
		pr.WriteString(strings.Repeat(".", pr.Len()))
		pr.WriteString(promptContinue)
		pr.WriteString(strings.Repeat(indent, max(0, depth)))
		return pr.String()
	}

	pr.WriteString(promptDefault)

	return pr.String()
}

func (cl *contLiner) promptString(str string) string {
	return cl.promptStringWithDepth(str, cl.depth)
}

func (cl *contLiner) Prompt(str string) (string, error) {
	line, err := cl.State.Prompt(cl.promptString(str))

	switch err {
	case nil:
		if cl.buffer != "" {
			cl.buffer = cl.buffer + "\n" + line
		} else {
			cl.buffer = line
		}
	case io.EOF:
		// when ^D
		if cl.buffer != "" {
			// cancel line continuation if in continuation
			cl.Accepted()
		}
		fmt.Println()
		fmt.Println("See you!")
		// else do nothing
	case liner.ErrPromptAborted:
		if cl.buffer != "" {
			cl.Accepted()
		} else {
			fmt.Println("(^D to quit)")
		}
	}

	return cl.buffer, err
}

func (cl *contLiner) Accepted() {
	cl.State.AppendHistory(strings.ReplaceAll(cl.buffer, "\n", " "))
	cl.Clear()
}

func (cl *contLiner) Clear() {
	cl.buffer = ""
	cl.depth = 0
	cl.prompt = ""
}

func (cl *contLiner) Close() {
	cl.Clear()
	cl.State.Close()
}

var errUnmatchedBraces = errors.New("unmatched braces")

func (cl *contLiner) Reindent(str string) error {
	oldDepth := cl.depth
	cl.depth = cl.CountDepth()

	if cl.depth < 0 {
		return errUnmatchedBraces
	}

	lines := strings.Split(cl.buffer, "\n")
	if len(lines) > 1 {
		lastLine := lines[len(lines)-1]
		if cl.depth < oldDepth {
			cursorUp()
			fmt.Printf("\r%s%s", cl.promptString(str), lastLine)
			eraseInLine()
			fmt.Println()
		} else if hasPrefix(lastLine, token.RPAREN, token.RBRACE) {
			cursorUp()
			fmt.Printf("\r%s%s", cl.promptStringWithDepth(str, cl.depth-1), lastLine)
			eraseInLine()
			fmt.Println()
		}
	}

	return nil
}

func (cl *contLiner) CountDepth() int {
	l := lexer.New(cl.buffer)
	depth := 0
	for {
		switch l.NextToken().Type {
		case token.LPAREN, token.LBRACE:
			depth++
		case token.RPAREN, token.RBRACE:
			depth--
		case token.EOF:
			return depth
		}
	}
}
