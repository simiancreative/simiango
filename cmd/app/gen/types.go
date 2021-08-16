package gen

type Template struct {
	Path    string
	Content string
}

type RequiredVar struct {
	Type    string
	Name    string
	Message string
}

type Erator struct {
	Name         string
	Desc         string
	Templates    []Template
	RequiredVars []RequiredVar
}

type Values map[string]interface{}
