package snippet

import (
	"fmt"
	"runtime"
	"strings"
)

const powershellSnippet = `using namespace System.Management.Automation
using namespace System.Management.Automation.Language

$_carapace_magick_completer = {
    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute("PSAvoidUsingInvokeExpression", "", Scope="Function", Target="*")]
    param($wordToComplete, $commandAst, $cursorPosition)
    $commandElements = $commandAst.CommandElements

    $elems = @()
    foreach ($_ in $commandElements) {
      if ($_.Extent.StartOffset -gt $cursorPosition) {
          break
      }
      $t = $_.Extent.Text
      if ($_.Extent.EndOffset -gt $cursorPosition) {
          $t = $t.Substring(0, $_.Extent.Text.get_Length() - ($_.Extent.EndOffset - $cursorPosition))
      }

      if ($t.Substring(0,1) -eq "'"){
        $t = $t.Substring(1)
      }
      if ($t.get_Length() -gt 0 -and $t.Substring($t.get_Length()-1) -eq "'"){
        $t = $t.Substring(0,$t.get_Length()-1)
      }
      if ($t.get_Length() -eq 0){
        $t = '""'
      }
      $elems += $t.replace('` + "`" + `,', ',') # quick fix
    }

    $completions = @(
      if (!$wordToComplete) {
        %[1]v ($elems[0] -replace ('\.exe$', '')) _carapace powershell $($elems| ForEach-Object {$_}) '' | ConvertFrom-Json | ForEach-Object { [CompletionResult]::new($_.CompletionText, $_.ListItemText.replace('` + "`" + `e[', "` + "`" + `e["), [CompletionResultType]::ParameterValue, $_.ToolTip.replace('` + "`" + `e[', "` + "`" + `e[")) }
      } else {
        %[1]v ($elems[0] -replace ('\.exe$', '')) _carapace powershell $($elems| ForEach-Object {$_}) | ConvertFrom-Json | ForEach-Object { [CompletionResult]::new($_.CompletionText, $_.ListItemText.replace('` + "`" + `e[', "` + "`" + `e["), [CompletionResultType]::ParameterValue, $_.ToolTip.replace('` + "`" + `e[', "` + "`" + `e[")) }
      }
    )

    if ($completions.count -eq 0) {
      return "" # prevent default file completion
    }

    $completions
}

%[2]v
`

func powershell() string {
	complete := make([]string, len(completerNames))
	for i, name := range completerNames {
		prefix := " # "
		if runtime.GOOS == "windows" {
			prefix = ""
		}
		complete[i] = fmt.Sprintf(`Register-ArgumentCompleter -Native -ScriptBlock $_carapace_magick_completer -CommandName '%v'%v'%v.exe'`, name, prefix, name)
	}
	return fmt.Sprintf(powershellSnippet, executable(), strings.Join(complete, "\n"))
}
