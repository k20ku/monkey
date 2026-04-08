package parser

import (
	"fmt"
	"strings"
	"time"
)

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) (string, time.Time) {
	incIdent()
	tracePrint("BEGIN " + msg)
	before := time.Now()
	return msg, before
}

func untrace(msg string, before time.Time) {
	after := time.Now()
	tracePrint(
		"END " + msg + fmt.Sprintf(" (%d μs)", (after.Sub(before)).Microseconds()),
	)
	decIdent()
}
