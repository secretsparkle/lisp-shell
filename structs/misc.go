package structs

type Control struct {
	Args []string
	Body []string
}

type Function struct {
	Name     string
	Args     []string
	Bindings map[string]string
	Body     List
}

func Maps() (map[string]rune, map[string]Function, map[string]string) {
	symbols := map[string]rune{
		"and":       'f',
		"car":       'f',
		"cdr":       'f',
		"cond":      'f',
		"cons":      'f',
		"defun":     'f',
		"defvar":    'f',
		"equal":     'f',
		"first":     'f',
		"if":        'f',
		"interpret": 'f',
		"last":      'f',
		"map":       'f',
		"quote":     'f',
		"rest":      'f',
		"reverse":   'f',
		"=":         'f',
		"+":         'f',
		"-":         'f',
		"*":         'f',
		"/":         'f',
		"<":         'f',
		">":         'f',
	}
	functions := make(map[string]Function)
	bindings := make(map[string]string)
	return symbols, functions, bindings
}
