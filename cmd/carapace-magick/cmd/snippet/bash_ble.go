package snippet

import (
	"fmt"
	"strings"
)

func bashBle() string {
	return fmt.Sprintf(`%v

_%[2]v_completion_ble() {
  if [[ ${BLE_ATTACHED-} ]]; then
    [[ :$comp_type: == *:auto:* ]] && return

    compopt -o ble/no-default
    bleopt complete_menu_style=desc

    local command="${COMP_WORDS[0]}"
    local compline="${COMP_LINE:0:${COMP_POINT}}"
    local IFS=$'\n'
    local c
    mapfile -t c < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %[2]v "${command}" _carapace bash-ble)
    [[ "${c[*]}" == "" ]] && c=() # fix for mapfile creating a non-empty array from empty command output

    local cand
    for cand in "${c[@]}"; do
      [ ! -z "$cand" ] && ble/complete/cand/yield mandb "${cand%%$'\t'*}" "${cand##*$'\t'}"
    done
  else
    complete -F _carapace_magick_completer %[3]v
  fi
}

complete -F _%[2]v_completion_ble %[3]v
`, bash(), executable(), strings.Join(completerNames, " "))
}
