package core

type Principal struct {
	Subject string
	Owner   string
	Name    string
	Email   string
	Scope   []string
	Roles   []string
}
