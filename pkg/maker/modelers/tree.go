package modelers

type Tree struct {
}

type ElementNode struct {
	Name    string `json:"name,omitempty"` // 名称，同一个应用中唯一
	element *Element
}

func (this_ *ElementNode) GetName() string {
	if this_ == nil {
		return ""
	}
	return this_.Name
}

func (this_ *ElementNode) SetName(name string) {
	if this_ == nil {
		return
	}
	this_.Name = name
}

func (this_ *ElementNode) GetElement() *Element {
	if this_ == nil {
		return nil
	}
	return this_.element
}

func (this_ *ElementNode) SetElement(element *Element) {
	if this_ == nil {
		return
	}
	this_.element = element
}

type ElementIFace interface {
	GetName() string
	SetName(name string)
	GetElement() *Element
	SetElement(element *Element)
}

type Element struct {
	Key      string     `json:"key,omitempty"`
	Text     string     `json:"text,omitempty"`
	IsType   bool       `json:"isType,omitempty"`
	IsPack   bool       `json:"isPack,omitempty"`
	IsModel  bool       `json:"isModel,omitempty"`
	Children []*Element `json:"children,omitempty"`
	Pack     *Pack      `json:"pack,omitempty"`
	parent   *Element
	Model    any `json:"model,omitempty"`
}

func (this_ *Element) GetParent() *Element {
	return this_.parent
}

func (this_ *Element) SetParent(parent *Element) {
	this_.parent = parent
}

type Pack struct {
	Name    string `json:"name,omitempty"`    // 名称，同一个应用中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
}
