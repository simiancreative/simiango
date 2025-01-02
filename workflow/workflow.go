package workflow

type Model struct {
	Name        string
	Description string
	Actions
}

type Actions map[string]Action

type Action struct {
	Args   ArgsList
	Runner func(args Args) (interface{}, error)
}
