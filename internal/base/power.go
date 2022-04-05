package base

import "github.com/gin-gonic/gin"

type ApiWorker struct {
	Apis    []string
	Power   *PowerAction
	Do      func(request *RequestBean, c *gin.Context) (res interface{}, err error)
	DoOther func(request *RequestBean, c *gin.Context)
}

type PowerAction struct {
	Action      string `json:"action,omitempty"`
	Text        string `json:"text,omitempty"`
	ShouldLogin bool   `json:"shouldLogin,omitempty"`
	StandAlone  bool   `json:"standAlone,omitempty"` // 单机是否可用
	Parent      *PowerAction
}

var (
	powers []*PowerAction
)

func addPower(power *PowerAction) *PowerAction {
	powers = append(powers, power)
	return power
}

func AppendPower(power *PowerAction) *PowerAction {
	return addPower(power)
}
func GetPowers() (ps []*PowerAction) {

	ps = append(ps, powers...)

	return
}
