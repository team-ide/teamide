package form

type Form struct {
	Fields []*Field `json:"fields,omitempty"`
}

type Field struct {
	Label        string      `json:"label,omitempty"`
	Name         string      `json:"name,omitempty"`
	BindName     string      `json:"bindName,omitempty"`
	Placeholder  string      `json:"placeholder,omitempty"`
	Type         string      `json:"type,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	IsNumber     bool        `json:"isNumber,omitempty"`
	Rules        []*Rule     `json:"rules,omitempty"`
	Options      []*Option   `json:"options,omitempty"`
}

type Rule struct {
	Required  bool   `json:"required,omitempty"`
	Length    int    `json:"length,omitempty"`
	Min       int    `json:"min,omitempty"`
	Max       int    `json:"max,omitempty"`
	MinLength int    `json:"minLength,omitempty"`
	MaxLength int    `json:"maxLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	Message   string `json:"message,omitempty"`
}

type Option struct {
	Text  string `json:"text,omitempty"`
	Value string `json:"value,omitempty"`
}
