package repl

import (
	"fmt"
	"io"
	"maps"
	"monkey/lexer"
	"monkey/token"
	"os"
	"path/filepath"
	"strings"

	// A well-known OSS golang repl "gore" uses "liner".
	// See [gore](https://github.com/x-motemen/gore/blob/main/liner.go#L11).
	liner "github.com/peterh/liner" // go get github.com/peterh/liner
)

const PROMPT = "monkey> "

var (
	history_filepath = filepath.Join(os.TempDir(), ".monkey_history")
	keywords         = maps.Keys(token.Keywords)
)

func Start() {

	// init/close liner
	line := liner.NewLiner()
	defer func() {
		// write all histories in this session.
		if hfd, err := os.Create(history_filepath); err == nil { // truncate history_file if it exists
			line.WriteHistory(hfd)
			hfd.Close()
		} else {
			fmt.Println("Error writing history file:\n", err)
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

	fmt.Println("Ctrl+D to exit!")

	// start repl
	for {
		if codeline, err := line.Prompt(PROMPT); err == nil {
			// lexer
			// give the lexer the code which the user typed
			l := lexer.New(codeline)
			// print all tokens
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				fmt.Printf("%+v\n", tok)
			}
			line.AppendHistory(codeline)

		} else if err == io.EOF {
			fmt.Println("(^D)")
			return

		} else if err == liner.ErrPromptAborted { // resume repl if aborted
			continue

		} else {
			fmt.Println("Error reading line: ", err)
		}

	}
}
