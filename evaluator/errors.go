package evaluator

import "fmt"

var (
	err_wrong_number_of_arguments = "wrong number of arguments. got=%d, want=%d" // got=%d want=%d
)

func serrorf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
