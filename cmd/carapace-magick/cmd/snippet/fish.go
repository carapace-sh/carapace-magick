package snippet

import (
	"fmt"
	"strings"
)

func fish() string {
	complete := make([]string, 0, len(completerNames)*2)
	for _, name := range completerNames {
		complete = append(complete,
			fmt.Sprintf(`complete -e %q`, name),
			fmt.Sprintf(`complete -c %q -f -a '(_carapace_magick_completer %q)' -r`, name, name),
		)
	}
	return fmt.Sprintf(`function _carapace_magick_completer
  set --local data
  IFS='' set data (echo (commandline -cp)'' | sed "s/ \$/ ''/" | xargs %[1]v $argv[1] _carapace fish 2>/dev/null)
  if [ $status -eq 1 ]
    IFS='' set data (echo (commandline -cp)"'" | sed "s/ \$/ ''/" | xargs %[1]v $argv[1] _carapace fish 2>/dev/null)
    if [ $status -eq 1 ]
      IFS='' set data (echo (commandline -cp)'"' | sed "s/ \$/ ''/" | xargs %[1]v $argv[1] _carapace fish 2>/dev/null)
    end
  end
  echo $data
end

%v
`, executable(), strings.Join(complete, "\n"))
}

func fishSingle(command string) string {
	return fmt.Sprintf(`complete -e %[2]q
complete -c %[2]q -f -a '(echo (commandline -cp)\'\' | sed "s/ \$/ \'\'/" | xargs %[1]v %[2]v _carapace fish 2>/dev/null)' -r
`, executable(), command)
}
