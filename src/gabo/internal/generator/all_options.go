package generator

// Option represents a generation option
type Option string

const (
	_LintDocker      Option = "lint-docker"
	_LintGo          Option = "lint-go"
	_LintMarkdown    Option = "lint-markdown"
	_LintPython      Option = "lint-python"
	_LintShellScript Option = "lint-shell-script"
	_LintSolidity    Option = "lint-solidity"
	_LintYaml        Option = "lint-yaml"

	_BuildDocker Option = "build-docker"

	// Make this code correct first
	// _LintAndroid  Option   = "lint-android"
	// _LintHtml     Option   = "lint-html"

	// TODO(ashishb): Enable these
	// _BuildAndroid Option = "build-android"
	// _BuildGo Option = "build-go"
	// _BuildPython  Option = "build-python"
	// _AutoTranslateAndroid Option = "auto-translate-android"
	// _TestGo Option = "test-go"
	// _TestPython Option = "test-python"
)

var _options = []Option{
	_LintDocker,
	_LintGo,
	_LintMarkdown,
	_LintPython,
	_LintShellScript,
	_LintSolidity,
	_LintYaml,

	_BuildDocker,
}

func GetOptions() []string {
	result := make([]string, 0, len(_options))
	for _, option := range _options {
		result = append(result, string(option))
	}
	return result
}

func IsValid(val string) bool {
	for _, option := range _options {
		if val == string(option) {
			return true
		}
	}
	return false
}
