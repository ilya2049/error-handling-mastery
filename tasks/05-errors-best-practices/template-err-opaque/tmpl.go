package tmpl

import (
	"errors"
	"regexp"
	"strings"
	"text/template"
)

var notDefinedRegexp = regexp.MustCompile(`function \".*\" not defined`)

// IsFunctionNotDefinedError говорит, является ли err ошибкой неопределённой в шаблоне функции.
func IsFunctionNotDefinedError(err error) bool {
	if err == nil {
		return false
	}

	return notDefinedRegexp.MatchString(err.Error())
}

// IsExecUnexportedFieldError говорит, является ли err template.ExecError,
// а именно ошибкой использования неэкспортируемого поля структуры.
func IsExecUnexportedFieldError(err error) bool {
	if err == nil {
		return false
	}

	var (
		execError template.ExecError
		isExec    bool
		isContain bool
	)

	isExec = errors.As(err, &execError)
	isContain = strings.Contains(err.Error(), "is an unexported field of struct")

	return isExec && isContain
}
