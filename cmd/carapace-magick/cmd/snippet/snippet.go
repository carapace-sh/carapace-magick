package snippet

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/carapace-sh/carapace/pkg/ps"
)

var completerNames = []string{"magick", "identify", "mogrify", "compare", "composite", "montage"}

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
	if shell == "" {
		shell = ps.DetermineShell()
	}
	switch shell {
	case "bash":
		return bash()
	case "bash-ble":
		return bashBle()
	case "zsh":
		return zsh()
	case "fish":
		return fish()
	case "elvish":
		return elvish()
	case "nushell":
		return nushell()
	case "powershell":
		return powershell()
	case "xonsh":
		return xonsh()
	case "oil":
		return oil()
	case "tcsh":
		return tcsh()
	case "export":
		return ""
	default:
		supported := []string{"bash", "bash-ble", "elvish", "fish", "nushell", "oil", "powershell", "tcsh", "xonsh", "zsh"}
		sort.Strings(supported)
		return fmt.Sprintf("# unsupported shell: '%v' [expected one of '%v']", shell, strings.Join(supported, "', '"))
	}
}
