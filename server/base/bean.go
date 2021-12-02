package base

type OBean struct {
	Text    string      `json:"text" column:"text"`
	Value   interface{} `json:"value" column:"value"`
	Comment string      `json:"comment" column:"comment"`
	Color   string      `json:"color" column:"color"`
}

func NewOBean(text string, value interface{}) (res OBean) {
	res = OBean{
		Text:  text,
		Value: value,
	}
	return
}

type UserTotalBean struct {
	User       *UserEntity         `json:"user"`
	Password   *UserPasswordEntity `json:"password"`
	Persona    *UserPersonaBean    `json:"persona"`
	Enterprise *UserEnterpriseBean `json:"enterprise"`
}

type UserPersonaBean struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Sex    int8    `json:"sex"`
	Photo  string  `json:"photo"`
	Height float32 `json:"height"`
	Weight float32 `json:"weight"`
}

type UserEnterpriseBean struct {
	Name   string                  `json:"name"`
	Salary float32                 `json:"salary"`
	Orgs   []UserEnterpriseOrgBean `json:"orgs"`
}

type UserEnterpriseOrgBean struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Position string `json:"position"`
}
