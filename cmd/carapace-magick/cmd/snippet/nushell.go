package snippet

import "fmt"

func nushell() string {
	return fmt.Sprintf(`let carapace_magick_completer = {|spans|
    %v $spans.0 _carapace nushell ...$spans | from json
}

mut current = (($env | default {} config).config | default {} completions)
$current.completions = ($current.completions | default {} external)
$current.completions.external = ($current.completions.external
|| default true enable
|# backwards compatible workaround for default, see nushell #15654
|| upsert completer { if $in == null { $carapace_magick_completer } else { $in } })

$env.config = $current
`, executable())
}

func nushellSingle(command string) string {
	return fmt.Sprintf(`let carapace_magick_completer = {|spans|
    %v %[2]v _carapace nushell ...$spans | from json
}

mut current = (($env | default {} config).config | default {} completions)
$current.completions = ($current.completions | default {} external)
$current.completions.external = ($current.completions.external
|| default true enable
|# backwards compatible workaround for default, see nushell #15654
|| upsert completer { if $in == null { $carapace_magick_completer } else { $in } })

$env.config = $current
`, executable(), command)
}
