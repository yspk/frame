package constant

import "strings"

const (
	UnknownError              = -1
	Success                   = 0
	Failure                   = 1
	StatusInternalServerError = 500
)

func TranslateErrCode(code int, extra ...string) string {
	var msg string
	switch code {
	case UnknownError:
		msg = "Unknown error"
	default:
	}

	if len(extra) > 0 {
		msg = msg + ": " + strings.Join(extra, ",")
	}
	return msg
}
