package maker

type Component struct {
	Fields  []*ComponentField  `json:"fields"`
	Methods []*ComponentMethod `json:"methods"`
}

type ComponentField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ComponentMethod struct {
	Name        string            `json:"name"`
	Args        []*ComponentField `json:"args"`
	ReturnTypes []string          `json:"returnTypes"`
	ThrowErrors []string          `json:"throwErrors"`
}

func (this_ *Component) ToContext() (ctx map[string]interface{}) {
	ctx = make(map[string]interface{})

	return ctx
}
