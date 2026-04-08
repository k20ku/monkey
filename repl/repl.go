package repl

import (
	"fmt"
	"io"
	"maps"
	"monkey/lexer"
	"monkey/parser"
	"monkey/token"
	"os"
	"path/filepath"
	"strings"

	// A well-known OSS golang repl "gore" uses "liner".
	// See [gore](https://github.com/x-motemen/gore/blob/main/liner.go#L11).
	liner "github.com/peterh/liner" // go get github.com/peterh/liner
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

type ReplMode string

const (
	LEX     ReplMode = "/lex"
	PARSE   ReplMode = "/parse"
	DEFAULT ReplMode = "/"
)

func (rm ReplMode) String() string {
	switch rm {
	case LEX:
		return "/lex"
	case PARSE:
		return "/parse"
	case DEFAULT:
		return "/"
	}
	return ""
}

var AllReplModes = []ReplMode{LEX, PARSE, DEFAULT}

var replMode ReplMode = DEFAULT

var PromptOn = map[ReplMode]string{
	LEX:     "monkey(LEX)> ",
	PARSE:   "monkey(PARSE)> ",
	DEFAULT: "monkey> ",
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

	// start repl
	for {
		codeline, err := line.Prompt(PromptOn[replMode])

		codeline = strings.TrimSpace(codeline)

		if changeMode(codeline) {
			mode := strings.Split(codeline, " ")[0]
			switch strings.ToLower(mode) {
			case LEX.String():
				replMode = LEX
			case PARSE.String():
				replMode = PARSE
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
			doOn(replMode, codeline)
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

func doOn(mode ReplMode, codeline string) {
	switch mode {
	case LEX:
		doOnLex(codeline)
	case PARSE:
		doOnParse(codeline)
	case DEFAULT:
		doOnDefault(codeline)
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

func doOnDefault(codeline string) {
	doOnParse(codeline)
}

func changeMode(codeline string) bool {
	return strings.HasPrefix(codeline, "/")
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
