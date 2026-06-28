package snippet

import (
	"fmt"
	"strings"
)

func bash() string {
	quoted := make([]string, len(completerNames))
	for i, name := range completerNames {
		quoted[i] = fmt.Sprintf("%q", name)
	}
	return fmt.Sprintf(`#!/bin/bash
_carapace_magick_completer() {
  export COMP_LINE
  export COMP_POINT
  export COMP_TYPE
  export COMP_WORDBREAKS

  local nospace data compline="${COMP_LINE:0:${COMP_POINT}}"
  local command="${COMP_WORDS[0]}"

  data=$(echo "${compline}''" | xargs %[1]v "${command}" _carapace bash 2>/dev/null)
  if [ $? -eq 1 ]; then
    data=$(echo "${compline}'" | xargs %[1]v "${command}" _carapace bash 2>/dev/null)
    if [ $? -eq 1 ]; then
    	data=$(echo "${compline}\"" | xargs %[1]v "${command}" _carapace bash 2>/dev/null)
    fi
  fi

  IFS=$'\001' read -r -d '' nospace data <<<"${data}"
  mapfile -t COMPREPLY < <(echo "${data}")
  unset COMPREPLY[-1]

  [ "${nospace}" = true ] && compopt -o nospace
  local IFS=$'\n'
  [[ "${COMPREPLY[*]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output
}

complete -o noquote -F _carapace_magick_completer %[2]v
`, executable(), strings.Join(quoted, " "))
}
