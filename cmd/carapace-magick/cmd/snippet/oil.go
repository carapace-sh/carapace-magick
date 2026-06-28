package snippet

import (
	"fmt"
	"strings"
)

func oil() string {
	quoted := make([]string, len(completerNames))
	for i, name := range completerNames {
		quoted[i] = fmt.Sprintf("%q", name)
	}
	return fmt.Sprintf(`#!/bin/osh
_carapace_magick_completer() {
  local command="${COMP_WORDS[0]}"
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %v "${command}" _carapace oil)
  [[ "${COMPREPLY[@]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output
  [[ ${COMPREPLY[0]} == *[/=@:.,$'\001'] ]] && compopt -o nospace
  # shellcheck disable=SC2206
  [[ ${#COMPREPLY[@]} -eq 1 ]] && COMPREPLY=(${COMPREPLY[@]%%$'\001'})
}

complete -F _carapace_magick_completer %v
`, executable(), strings.Join(quoted, " "))
}

func oilSingle(command string) string {
	return fmt.Sprintf(`#!/bin/osh
_%[2]v_completion() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %[1]v %[2]v _carapace oil)
  [[ "${COMPREPLY[@]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output
  [[ ${COMPREPLY[0]} == *[/=@:.,$'\001'] ]] && compopt -o nospace
  # shellcheck disable=SC2206
  [[ ${#COMPREPLY[@]} -eq 1 ]] && COMPREPLY=(${COMPREPLY[@]%%$'\001'})
}

complete -F _%[2]v_completion %[2]v
`, executable(), command)
}
