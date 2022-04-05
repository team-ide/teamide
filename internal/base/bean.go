package base

type OBean struct {
	Text    string      `json:"text" column:"text,omitempty"`
	Value   interface{} `json:"value" column:"value,omitempty"`
	Comment string      `json:"comment" column:"comment,omitempty"`
	Color   string      `json:"color" column:"color,omitempty"`
}

func NewOBean(text string, value interface{}) (res OBean) {
	res = OBean{
		Text:  text,
		Value: value,
	}
	return
}
