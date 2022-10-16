package common

type Code struct {
	Dir       string
	Name      string
	Namespace string
	Header    string
	Body      string
	Footer    string
	Tab       int
	Other     map[string]*Code
	CodeType  CodeType
}

type CodeType string

var (
	CodeTypeJava = CodeType("java")
	CodeTypeGo   = CodeType("go")
	CodeTypeXml  = CodeType("xml")
	CodeTypeYaml = CodeType("yaml")
	CodeTypeJs   = CodeType("js")
)

func (this_ *Code) GetOther(key string) (other *Code) {
	if this_.Other == nil {
		this_.Other = make(map[string]*Code)
	}
	other = this_.Other[key]
	return
}

func (this_ *Code) SetOther(key string, other *Code) {
	if this_.Other == nil {
		this_.Other = make(map[string]*Code)
	}
	this_.Other[key] = other
	return
}

func (this_ *Code) ToContent() (content string) {

	return
}
