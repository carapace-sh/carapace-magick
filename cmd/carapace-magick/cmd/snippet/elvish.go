package snippet

import (
	"fmt"
	"runtime"
	"strings"
)

func elvish() string {
	quoted := make([]string, len(completerNames))
	for i, name := range completerNames {
		quoted[i] = fmt.Sprintf("%q", name)
	}
	windowsSnippet := ""
	if runtime.GOOS == "windows" {
		windowsSnippet = "\n    set edit:completion:arg-completer[$c.exe] = $edit:completion:arg-completer[$c]\n"
	}
	return fmt.Sprintf(`put %[2]v | each {|c|
    set edit:completion:arg-completer[$c] = {|@arg|
        %[1]v $c _carapace elvish (all $arg) | from-json | each {|completion|
    		put $completion[Messages] | all (one) | each {|m|
    			edit:notify (styled "error: " red)$m
    		}
    		if (not-eq $completion[Usage] "") {
    			edit:notify (styled "usage: " $completion[DescriptionStyle])$completion[Usage]
    		}
    		put $completion[Candidates] | all (one) | peach {|c|
    			if (eq $c[Description] "") {
    		    	edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style]) &code-suffix=$c[CodeSuffix]
    			} else {
    		    	edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style])(styled " " $completion[DescriptionStyle]" bg-default")(styled "("$c[Description]")" $completion[DescriptionStyle]) &code-suffix=$c[CodeSuffix]
    			}
    		}
        }
    }%[3]v
}
`, executable(), strings.Join(quoted, " "), windowsSnippet)
}
