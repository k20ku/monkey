package repl

import (
	"fmt"
	"io"
	"maps"
	"slices"

	"os"
	"path/filepath"
	"strings"

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
	history_filepath = filepath.Join(os.TempDir(), ".monkey_history")
	keywords         = maps.Keys(token.Keywords)
)

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

type ReplMode struct {
	mode   string
	prompt string
}

var AllReplModes = []*ReplMode{}

func registerMode(mode string, prompt string) *ReplMode {
	m := &ReplMode{mode: mode, prompt: prompt}
	AllReplModes = append(AllReplModes, m)
	return m
}

var (
	LEX     = registerMode("/lex", "LEX")
	PARSE   = registerMode("/parse", "PARSE")
	EVAL    = registerMode("/eval", "EVAL")
	DEFAULT = registerMode("/", "")
)

func (rm *ReplMode) String() string {
	if slices.Contains(AllReplModes, rm) {
		return rm.mode
	}
	return ""
}

func (rm *ReplMode) Prompt() string {
	if rm == DEFAULT {
		return "monkey> "
	}

	return "monkey(" + rm.prompt + ")> "
}

func Start() {

	// init/close liner
	line := liner.NewLiner()
	defer func() {
		// write all histories in this session.
		if hfd, err := os.Create(history_filepath); err == nil { // truncate history_file if it exists
			line.WriteHistory(hfd)
			hfd.Close()
		} else {
			fmt.Println(" writing history file errors:\n\t", err.Error())
		}
		line.Close()
	}()

	// setting liner
	line.SetMultiLineMode(true)
	line.SetCtrlCAborts(true)
	line.SetCompleter(func(line string) (c []string) {
		for keyword := range keywords {
			if strings.HasPrefix(keyword, strings.ToLower(line)) {
				c = append(c, keyword)
			}
		}
		return
	})

	// read history
	if hfd, err := os.Open(history_filepath); err == nil {
		line.ReadHistory(hfd)
		hfd.Close()
	}

	// init mode
	var replMode *ReplMode = DEFAULT
	env := object.NewEnvironment()
	// start repl
	for {
		codeline, err := line.Prompt(replMode.Prompt())

		codeline = strings.TrimSpace(codeline)

		if changeMode(codeline) {
			mode := strings.Split(codeline, " ")[0]
			switch strings.ToLower(mode) {
			case LEX.String():
				replMode = LEX
			case PARSE.String():
				replMode = PARSE
			case EVAL.String():
				replMode = EVAL
			case DEFAULT.String():
				replMode = DEFAULT
			default:
				fmt.Printf("\t"+"No Mode for '%s'!"+"\n", mode)
				modes := []string{}
				for _, mode := range AllReplModes {
					modes = append(modes, "'"+mode.String()+"'")
				}
				fmt.Printf("\t"+"Only Either of %s is Allowed!"+"\n", strings.Join(modes, ", "))
			}

			continue
		}

		switch err {
		case nil:
			doOn(replMode, codeline, env)
			line.AppendHistory(codeline)

		case io.EOF:
			fmt.Println("(^D)")
			return

		case liner.ErrPromptAborted: // resume repl if aborted
			continue

		default:
			fmt.Println(" reading line errors:\n\t", err)
		}
	}
}

func doOn(mode *ReplMode, codeline string, env *object.Environment) {
	switch mode {
	case LEX:
		doOnLex(codeline)
	case PARSE:
		doOnParse(codeline)
	case EVAL:
		doOnEval(codeline, env)
	case DEFAULT:
		doOnDefault(codeline, env)
	}
}

func doOnLex(codeline string) {
	l := lexer.New(codeline)
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
}

func doOnParse(codeline string) {
	l := lexer.New(codeline)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	io.WriteString(os.Stdout, program.String())
	io.WriteString(os.Stdout, "\n")
}

func doOnEval(codeline string, env *object.Environment) {
	l := lexer.New(codeline)
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

func doOnDefault(codeline string, env *object.Environment) {
	doOnEval(codeline, env)
}

func changeMode(codeline string) bool {
	return strings.HasPrefix(codeline, "/")
}

func printParserErrors(out io.Writer, errors []string) {
	// io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
