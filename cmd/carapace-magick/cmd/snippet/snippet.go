package snippet

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/carapace-sh/carapace/pkg/ps"
)

var completerNames = []string{"carapace-magick", "magick", "identify", "mogrify", "compare", "composite", "montage"}

func executable() string {
	e, err := os.Executable()
	if err != nil {
		return "carapace-magick"
	}
	return filepath.Base(e)
}

// Snippet generates a multi-completer shell snippet that registers all
// ImageMagick command completers at once.
func Snippet(shell string) string {
	return snippetFor(shell, "", false)
}

// SingleSnippet generates a shell snippet that registers completion for
// a single ImageMagick command (e.g. `carapace-magick identify _carapace bash`).
func SingleSnippet(shell, command string) string {
	return snippetFor(shell, command, true)
}

func snippetFor(shell, command string, single bool) string {
	if shell == "" {
		shell = ps.DetermineShell()
	}
	switch shell {
	case "bash":
		if single {
			return bashSingle(command)
		}
		return bash()
	case "bash-ble":
		if single {
			return bashBleSingle(command)
		}
		return bashBle()
	case "zsh":
		if single {
			return zshSingle(command)
		}
		return zsh()
	case "fish":
		if single {
			return fishSingle(command)
		}
		return fish()
	case "elvish":
		if single {
			return elvishSingle(command)
		}
		return elvish()
	case "nushell":
		if single {
			return nushellSingle(command)
		}
		return nushell()
	case "powershell":
		if single {
			return powershellSingle(command)
		}
		return powershell()
	case "xonsh":
		if single {
			return xonshSingle(command)
		}
		return xonsh()
	case "oil":
		if single {
			return oilSingle(command)
		}
		return oil()
	case "tcsh":
		if single {
			return tcshSingle(command)
		}
		return tcsh()
	case "export":
		return ""
	default:
		supported := []string{"bash", "bash-ble", "elvish", "fish", "nushell", "oil", "powershell", "tcsh", "xonsh", "zsh"}
		sort.Strings(supported)
		return fmt.Sprintf("# unsupported shell: '%v' [expected one of '%v']", shell, strings.Join(supported, "', '"))
	}
}
