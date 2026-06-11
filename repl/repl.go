package repl

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"strings"

	"os"
	"path/filepath"

	// A well-known OSS golang repl "gore" uses "liner".
	// See [gore](https://github.com/x-motemen/gore/blob/main/liner.go#L11).
	liner "github.com/peterh/liner" // go get github.com/peterh/liner

	"github.com/k20ku/monkey/evaluator"
	"github.com/k20ku/monkey/lexer"
	"github.com/k20ku/monkey/object"
	"github.com/k20ku/monkey/parser"
	"github.com/k20ku/monkey/token"
)

var (
	keywords = maps.Keys(token.Keywords)
)

func Start() {
	// init/close liner
	cl := newContLiner()
	defer cl.Close()

	cl.SetCompleter(func(line string) []string {
		var c []string
		for keyword := range keywords {
			if strings.HasPrefix(keyword, strings.ToLower(line)) {
				c = append(c, keyword)
			}
		}
		return c
	})
	// read history
	var historyFile string
	home, err := homeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "REPL: home: %+v", err)
	} else {
		historyFile = filepath.Join(home, "history")
		f, err := os.Open(historyFile)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "REPL: %+v\n", err)
			}
		} else {
			_, err := cl.ReadHistory(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "REPL: while reading history: %+v", err)
			}
			f.Close()
		}
	}

	// init env
	env := object.NewEnvironment()
	// start repl
	for {
		in, err := cl.Prompt("")

		if err != nil {
			if err == io.EOF {
				fmt.Println("(^D)")
				break
			} else if err == liner.ErrPromptAborted {
				continue
			}
			fmt.Fprintf(os.Stderr, "REPL: %+v\n", err)
		}

		if in == "" {
			continue
		}

		if err := cl.Reindent(""); err != nil {
			cl.Clear()
			continue
		}

		if cl.CountDepth() < 0 {
			fmt.Printf("%s: %s\n", in, errUnmatchedBraces)
			cl.Clear()
			continue
		} else if cl.CountDepth() != 0 {
			continue
		}

		Eval(cl.buffer, env)

		cl.Accepted()
		if historyFile != "" {
			err := os.MkdirAll(filepath.Dir(historyFile), 0o755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
			} else {
				f, err := os.Create(historyFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
				} else {
					_, err := cl.WriteHistory(f)
					if err != nil {
						fmt.Fprintf(os.Stderr, "while saving history: %s", err)
					}
					f.Close()
				}
			}
		}
	}

}

func homeDir() (home string, err error) {
	home = os.Getenv("MONKEY_HOME")
	if home != "" {
		return
	}

	var baseDir string

	baseDir = os.Getenv("XDG_DATA_HOME")
	if baseDir != "" {
		home = filepath.Join(baseDir, "monkey")

		return
	}

	baseDir, err = os.UserHomeDir()
	if err != nil {
		return
	}

	home = filepath.Join(baseDir, ".monkey")
	return
}

func Lex(in string) {
	l := lexer.New(in)
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
}

func Parse(in string) {
	l := lexer.New(in)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	b, err := json.MarshalIndent(program, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Stdout.Write(b)
	io.WriteString(os.Stdout, "\n")
}

func Eval(in string, env *object.Environment) {
	l := lexer.New(in)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(os.Stdout, evaluated.Inspect())
		io.WriteString(os.Stdout, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	// io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
