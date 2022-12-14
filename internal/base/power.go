package base

import "github.com/gin-gonic/gin"

type ApiWorker struct {
	Power        *PowerAction
	Do           func(request *RequestBean, c *gin.Context) (res interface{}, err error)
	IsGet        bool
	IsWebSocket  bool
	IsUpload     bool
	NotRecodeLog bool `json:"notRecodeLog"`
}

type PowerAction struct {
	Action       string       `json:"action,omitempty"`
	Text         string       `json:"text,omitempty"`
	ShouldLogin  bool         `json:"shouldLogin,omitempty"`
	StandAlone   bool         `json:"standAlone,omitempty"` // 单机是否可用
	ShouldPower  bool         `json:"shouldPower,omitempty"`
	ParentAction string       `json:"parentAction,omitempty"`
	Parent       *PowerAction `json:"-"`
}

var (
	powers []*PowerAction
)

func addPower(power *PowerAction) *PowerAction {
	if power.Parent != nil {
		power.ParentAction = power.Parent.Action
	}
	powers = append(powers, power)
	return power
}

func AppendPower(power *PowerAction) *PowerAction {
	if power.Parent != nil {
		power.ParentAction = power.Parent.Action
		power.Action = power.Parent.Action + "/" + power.Action
	}
	powers = append(powers, power)
	return power
}
func GetPowers() (ps []*PowerAction) {

	ps = powers

	return
}
