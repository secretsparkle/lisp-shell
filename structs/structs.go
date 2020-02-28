package structs

type Control struct {
	Args []string
	Body []string
}

type Function struct {
	Name     string
	Args     []string
	Bindings map[string]string
	Body     []string
}
