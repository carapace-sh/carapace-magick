package snippet

import (
	"fmt"
	"strings"
)

func tcsh() string {
	lines := make([]string, len(completerNames))
	for i, name := range completerNames {
		lines[i] = fmt.Sprintf("complete \"%v\" 'p@*@`echo \"$COMMAND_LINE'\"''\"'\" | xargs %v \"%v\" _carapace tcsh `@@' ;", name, executable(), name)
	}
	return strings.Join(lines, "\n")
}

func tcshSingle(command string) string {
	return fmt.Sprintf("complete \"%[2]v\" 'p@*@`echo \"$COMMAND_LINE'\"''\"'\" | xargs %[1]v \"%[2]v\" _carapace tcsh `@@' ;", executable(), command)
}
